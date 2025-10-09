package services

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	modelsUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	modelsUserContactRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relations"
	modelsUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"

	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoContact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repoUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_category_relations"
	repoUserContactRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_contact_relations"
	repoUserFull "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_full_repositories"
)

type UserFullService interface {
	CreateFull(ctx context.Context, user *modelsUserFull.UserFull) (*modelsUserFull.UserFull, error)
}

type userFullService struct {
	repoUser           repoUserFull.UserFullRepository
	repoAddress        repoAddress.AddressRepository
	repoContact        repoContact.ContactRepository
	repoUserCatRel     repoUserCatRel.UserCategoryRelationRepository
	repoUserContactRel repoUserContactRel.UserContactRelationRepository
	hasher             auth.PasswordHasher
}

func NewUserFullService(
	repoUser repoUserFull.UserFullRepository,
	repoAddress repoAddress.AddressRepository,
	repoContact repoContact.ContactRepository,
	repoUserCatRel repoUserCatRel.UserCategoryRelationRepository,
	repoUserContactRel repoUserContactRel.UserContactRelationRepository,
	hasher auth.PasswordHasher,
) UserFullService {
	return &userFullService{
		repoUser:           repoUser,
		repoAddress:        repoAddress,
		repoContact:        repoContact,
		repoUserCatRel:     repoUserCatRel,
		repoUserContactRel: repoUserContactRel,
		hasher:             hasher,
	}
}

func (s *userFullService) CreateFull(ctx context.Context, userFull *modelsUserFull.UserFull) (*modelsUserFull.UserFull, error) {
	if userFull == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := userFull.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	// Hash da senha
	hashed, err := s.hasher.Hash(userFull.User.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}
	userFull.User.Password = hashed

	// Inicia transação
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

	// Criação do usuário
	createdUser, err := s.repoUser.CreateTx(ctx, tx, userFull.User)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Criação do endereço
	userFull.Address.UserID = utils.StrToPtr(createdUser.UID)
	if err := userFull.Address.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("endereço inválido: %w", err))
	}
	createdAddress, err := s.repoAddress.CreateTx(ctx, tx, userFull.Address)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Criação do contato
	if err := userFull.Contact.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("contato inválido: %w", err))
	}
	createdContact, err := s.repoContact.CreateTx(ctx, tx, userFull.Contact)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Relação user-contact
	relation := &modelsUserContactRel.UserContactRelations{
		UserID:    createdUser.UID,
		ContactID: createdContact.ID,
	}
	if err := relation.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("relação usuário-contato inválida: %w", err))
	}
	if _, err := s.repoUserContactRel.CreateTx(ctx, tx, relation); err != nil {
		return nil, commitOrRollback(err)
	}

	// Relações user-category
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

	// Commit final da transação
	if err := commitOrRollback(nil); err != nil {
		return nil, err // garante que o objeto não seja retornado se o commit falhar
	}

	return &modelsUserFull.UserFull{
		User:       createdUser,
		Address:    createdAddress,
		Contact:    createdContact,
		Categories: userFull.Categories,
	}, nil
}
