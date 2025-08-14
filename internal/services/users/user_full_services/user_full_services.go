package services

import (
	"context"
	"errors"
	"fmt"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	models_user_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	models_user_full "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	repo_address "github.com/WagaoCarvalho/backend_store_go/internal/repositories/address"
	repo_contact "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	repo_relation "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	repo_user_full "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_full_repositories"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type UserFullService interface {
	CreateFull(ctx context.Context, user *models_user_full.UserFull) (*models_user_full.UserFull, error)
}

type userFullService struct {
	repo_user         repo_user_full.UserFullRepository
	repo_address      repo_address.AddressRepository
	repo_contact      repo_contact.ContactRepository
	repo_user_cat_rel repo_relation.UserCategoryRelationRepository
	logger            *logger.LoggerAdapter
	hasher            auth.PasswordHasher
}

func NewUserFullService(
	repo_user repo_user_full.UserFullRepository,
	repo_address repo_address.AddressRepository,
	repo_contact repo_contact.ContactRepository,
	repo_user_cat_rel repo_relation.UserCategoryRelationRepository,
	logger *logger.LoggerAdapter,
	hasher auth.PasswordHasher,
) UserFullService {
	return &userFullService{
		repo_user:         repo_user,
		repo_address:      repo_address,
		repo_contact:      repo_contact,
		repo_user_cat_rel: repo_user_cat_rel,
		logger:            logger,
		hasher:            hasher,
	}
}

func (s *userFullService) CreateFull(ctx context.Context, user_full *models_user_full.UserFull) (*models_user_full.UserFull, error) {
	ref := "[userService - CreateFull] - "

	logFields := map[string]any{}
	if user_full != nil && user_full.User != nil {
		logFields["username"] = user_full.User.Username
		logFields["email"] = user_full.User.Email
	}

	s.logger.Info(ctx, ref+logger.LogCreateInit, logFields)

	if err := user_full.Validate(); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogValidateError, logFields)
		return nil, err
	}

	if user_full.User.Password != "" {
		hashed, err := s.hasher.Hash(user_full.User.Password)
		if err != nil {
			s.logger.Error(ctx, err, ref+logger.LogPasswordInvalid, map[string]any{
				"email": user_full.User.Email,
			})
			return nil, fmt.Errorf("erro ao hashear senha: %w", err)
		}
		user_full.User.Password = hashed
	}

	tx, err := s.repo_user.BeginTx(ctx)
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

	createdUser, err := s.repo_user.CreateTx(ctx, tx, user_full.User)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	user_full.Address.UserID = utils.ToPointer(createdUser.UID)

	if err := user_full.Address.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("endereço inválido: %w", err))
	}

	createdAddress, err := s.repo_address.CreateTx(ctx, tx, user_full.Address)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	user_full.Contact.UserID = utils.ToPointer(createdUser.UID)

	if err := user_full.Contact.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("contato inválido: %w", err))
	}

	createdContact, err := s.repo_contact.CreateTx(ctx, tx, user_full.Contact)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	for _, category := range user_full.Categories {
		relation := &models_user_cat_rel.UserCategoryRelations{
			UserID:     createdUser.UID,
			CategoryID: int64(category.ID),
		}

		if err := relation.Validate(); err != nil {
			return nil, commitOrRollback(fmt.Errorf("relação usuário-categoria inválida: %w", err))
		}

		if _, err := s.repo_user_cat_rel.CreateTx(ctx, tx, relation); err != nil {
			return nil, commitOrRollback(err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdUser.UID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
	})

	result := &models_user_full.UserFull{
		User:       createdUser,
		Address:    createdAddress,
		Contact:    createdContact,
		Categories: user_full.Categories,
	}

	return result, commitOrRollback(nil)
}
