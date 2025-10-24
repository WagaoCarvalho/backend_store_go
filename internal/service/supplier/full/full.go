package services

import (
	"context"
	"errors"
	"fmt"

	modelsCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	modelsFull "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoContact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repoRelation "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/category_relation"
	repoContactRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/contact_relation"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/full"
)

type SupplierFull interface {
	CreateFull(ctx context.Context, supplierFull *modelsFull.SupplierFull) (*modelsFull.SupplierFull, error)
}

type supplierFull struct {
	repoSupplier   repoSupplier.SupplierFull
	repoAddress    repoAddress.Address
	repoContact    repoContact.Contact
	repoCatRel     repoRelation.SupplierCategoryRelation
	repoContactRel repoContactRel.SupplierContactRelation
}

func NewSupplierFull(
	repoSupplier repoSupplier.SupplierFull,
	repoAddress repoAddress.Address,
	repoContact repoContact.Contact,
	repoCatRel repoRelation.SupplierCategoryRelation,
	repoContactRel repoContactRel.SupplierContactRelation,
) SupplierFull {
	return &supplierFull{
		repoSupplier:   repoSupplier,
		repoAddress:    repoAddress,
		repoContact:    repoContact,
		repoCatRel:     repoCatRel,
		repoContactRel: repoContactRel,
	}
}

func (s *supplierFull) CreateFull(ctx context.Context, supplierFull *modelsFull.SupplierFull) (*modelsFull.SupplierFull, error) {
	if supplierFull == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := supplierFull.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	// Inicia transação
	tx, err := s.repoSupplier.BeginTx(ctx)
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

	// Criação do fornecedor
	createdSupplier, err := s.repoSupplier.CreateTx(ctx, tx, supplierFull.Supplier)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Criação do endereço
	supplierFull.Address.SupplierID = utils.StrToPtr(createdSupplier.ID)
	if err := supplierFull.Address.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("endereço inválido: %w", err))
	}
	createdAddress, err := s.repoAddress.CreateTx(ctx, tx, supplierFull.Address)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Criação do contato
	if err := supplierFull.Contact.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("contato inválido: %w", err))
	}
	createdContact, err := s.repoContact.CreateTx(ctx, tx, supplierFull.Contact)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	// Relações fornecedor-categoria
	for _, category := range supplierFull.Categories {
		relation := &modelsCatRel.SupplierCategoryRelation{
			SupplierID: createdSupplier.ID,
			CategoryID: int64(category.ID),
		}

		if err := relation.Validate(); err != nil {
			return nil, commitOrRollback(fmt.Errorf("relação fornecedor-categoria inválida: %w", err))
		}

		if _, err := s.repoCatRel.CreateTx(ctx, tx, relation); err != nil {
			return nil, commitOrRollback(err)
		}
	}

	// Commit final
	if err := commitOrRollback(nil); err != nil {
		return nil, err // garante que não retorne o objeto se o commit falhar
	}

	return &modelsFull.SupplierFull{
		Supplier:   createdSupplier,
		Address:    createdAddress,
		Contact:    createdContact,
		Categories: supplierFull.Categories,
	}, nil
}
