package services

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	modelsUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	modelsUserContactRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	modelsUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/full"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"

	repoAddressTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/address"
	repoContactTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/contact"
	repoUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/category_relation"
	repoUserContactRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/contact_relation"
	repoUserFull "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/full"
)

type UserFull interface {
	CreateFull(ctx context.Context, user *modelsUserFull.UserFull) (*modelsUserFull.UserFull, error)
}

type userFull struct {
	repoUser           repoUserFull.UserFull
	repoAddressTx      repoAddressTx.AddressTx
	repoContactTx      repoContactTx.ContactTx
	repoUserCatRel     repoUserCatRel.UserCategoryRelation
	repoUserContactRel repoUserContactRel.UserContactRelation
	hasher             auth.PasswordHasher
}

func NewUserFull(
	repoUser repoUserFull.UserFull,
	repoAddressTx repoAddressTx.AddressTx,
	repoContactTx repoContactTx.ContactTx,
	repoUserCatRel repoUserCatRel.UserCategoryRelation,
	repoUserContactRel repoUserContactRel.UserContactRelation,
	hasher auth.PasswordHasher,
) UserFull {
	return &userFull{
		repoUser:           repoUser,
		repoAddressTx:      repoAddressTx,
		repoContactTx:      repoContactTx,
		repoUserCatRel:     repoUserCatRel,
		repoUserContactRel: repoUserContactRel,
		hasher:             hasher,
	}
}

func (s *userFull) CreateFull(ctx context.Context, userFull *modelsUserFull.UserFull) (*modelsUserFull.UserFull, error) {
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
	createdAddress, err := s.repoAddressTx.CreateTx(ctx, tx, userFull.Address)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Criação do contato
	if err := userFull.Contact.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("contato inválido: %w", err))
	}
	createdContact, err := s.repoContactTx.CreateTx(ctx, tx, userFull.Contact)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Relação user-contact
	relation := &modelsUserContactRel.UserContactRelation{
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
		relation := &modelsUserCatRel.UserCategoryRelation{
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
