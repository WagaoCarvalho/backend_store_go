package services

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	modelsUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	modelsUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
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
	hasher         auth.PasswordHasher
}

func NewUserFullService(
	repoUser repoUserFull.UserFullRepository,
	repoAddress repoAddress.AddressRepository,
	repoContact repoContact.ContactRepository,
	repoUserCatRel repoRelation.UserCategoryRelationRepository,
	hasher auth.PasswordHasher,
) UserFullService {
	return &userFullService{
		repoUser:       repoUser,
		repoAddress:    repoAddress,
		repoContact:    repoContact,
		repoUserCatRel: repoUserCatRel,
		hasher:         hasher,
	}
}

func (s *userFullService) CreateFull(ctx context.Context, userFull *modelsUserFull.UserFull) (*modelsUserFull.UserFull, error) {

	if userFull == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := userFull.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	hashed, err := s.hasher.Hash(userFull.User.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}
	userFull.User.Password = hashed

	tx, err := s.repoUser.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	if tx == nil {
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
				return fmt.Errorf("%v; rollback error: %w", err, rbErr)
			}
			return err
		}
		if cErr := tx.Commit(ctx); cErr != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
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

	result := &modelsUserFull.UserFull{
		User:       createdUser,
		Address:    createdAddress,
		Contact:    createdContact,
		Categories: userFull.Categories,
	}

	return result, commitOrRollback(nil)
}
