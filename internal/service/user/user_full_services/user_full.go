package services

import (
	"context"
	"errors"
	"fmt"

	modelsUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	modelsUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoContact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repoRelation "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_category_relations"
	repoUserFull "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_full_repositories"
)

type UserFullService interface {
	CreateFull(ctx context.Context, user *modelsUserFull.UserFull) (*modelsUserFull.UserFull, error)
}

type userFullService struct {
	repoUser       repoUserFull.UserFullRepository
	repoAddress    repoAddress.AddressRepository
	repoContact    repoContact.ContactRepository
	repoUserCatRel repoRelation.UserCategoryRelationRepository
	logger         *logger.LogAdapter
	hasher         auth.PasswordHasher
}

func NewUserFullService(
	repoUser repoUserFull.UserFullRepository,
	repoAddress repoAddress.AddressRepository,
	repoContact repoContact.ContactRepository,
	repoUserCatRel repoRelation.UserCategoryRelationRepository,
	logger *logger.LogAdapter,
	hasher auth.PasswordHasher,
) UserFullService {
	return &userFullService{
		repoUser:       repoUser,
		repoAddress:    repoAddress,
		repoContact:    repoContact,
		repoUserCatRel: repoUserCatRel,
		logger:         logger,
		hasher:         hasher,
	}
}

func (s *userFullService) CreateFull(ctx context.Context, userFull *modelsUserFull.UserFull) (*modelsUserFull.UserFull, error) {
	ref := "[userService - CreateFull] - "

	logFields := map[string]any{}
	if userFull != nil && userFull.User != nil {
		logFields["username"] = userFull.User.Username
		logFields["email"] = userFull.User.Email
	}

	s.logger.Info(ctx, ref+logger.LogCreateInit, logFields)

	if err := userFull.Validate(); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogValidateError, logFields)
		return nil, err
	}

	if userFull.User.Password != "" {
		hashed, err := s.hasher.Hash(userFull.User.Password)
		if err != nil {
			s.logger.Error(ctx, err, ref+logger.LogPasswordInvalid, map[string]any{
				"email": userFull.User.Email,
			})
			return nil, fmt.Errorf("erro ao hashear senha: %w", err)
		}
		userFull.User.Password = hashed
	}

	tx, err := s.repoUser.BeginTx(ctx)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogTransactionInitError, nil)
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	if tx == nil {
		s.logger.Error(ctx, errors.New("transação nula"), ref+logger.LogTransactionNull, nil)
		return nil, errors.New("transação inválida")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	commitOrRollback := func(err error) error {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Error(ctx, rbErr, ref+logger.LogRollbackError, nil)
				return fmt.Errorf("%v; rollback error: %w", err, rbErr)
			}
			return err
		}
		if cErr := tx.Commit(ctx); cErr != nil {
			s.logger.Error(ctx, cErr, ref+logger.LogCommitError, nil)

			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Error(ctx, rbErr, ref+logger.LogRollbackErrorAfterCommitFail, nil)
				return fmt.Errorf("erro ao commitar transação: %v; rollback error: %w", cErr, rbErr)
			}
			return fmt.Errorf("erro ao commitar transação: %w", cErr)
		}
		return nil
	}

	createdUser, err := s.repoUser.CreateTx(ctx, tx, userFull.User)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	userFull.Address.UserID = utils.StrToPtr(createdUser.UID)

	if err := userFull.Address.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("endereço inválido: %w", err))
	}

	createdAddress, err := s.repoAddress.CreateTx(ctx, tx, userFull.Address)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	userFull.Contact.UserID = utils.StrToPtr(createdUser.UID)

	if err := userFull.Contact.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("contato inválido: %w", err))
	}

	createdContact, err := s.repoContact.CreateTx(ctx, tx, userFull.Contact)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	for _, category := range userFull.Categories {
		relation := &modelsUserCatRel.UserCategoryRelations{
			UserID:     createdUser.UID,
			CategoryID: int64(category.ID),
		}

		if err := relation.Validate(); err != nil {
			return nil, commitOrRollback(fmt.Errorf("relação usuário-categoria inválida: %w", err))
		}

		if _, err := s.repoUserCatRel.CreateTx(ctx, tx, relation); err != nil {
			return nil, commitOrRollback(err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdUser.UID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
	})

	result := &modelsUserFull.UserFull{
		User:       createdUser,
		Address:    createdAddress,
		Contact:    createdContact,
		Categories: userFull.Categories,
	}

	return result, commitOrRollback(nil)
}
