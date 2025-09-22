package services

import (
	"context"
	"errors"
	"fmt"

	modelsCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	modelsFull "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoContact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repoRelation "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category_relations"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_full_repositories"
)

type SupplierFullService interface {
	CreateFull(ctx context.Context, supplierFull *modelsFull.SupplierFull) (*modelsFull.SupplierFull, error)
}

type supplierFullService struct {
	repoSupplier repoSupplier.SupplierFullRepository
	repoAddress  repoAddress.AddressRepository
	repoContact  repoContact.ContactRepository
	repoCatRel   repoRelation.SupplierCategoryRelationRepository
}

func NewSupplierFullService(
	repoSupplier repoSupplier.SupplierFullRepository,
	repoAddress repoAddress.AddressRepository,
	repoContact repoContact.ContactRepository,
	repoCatRel repoRelation.SupplierCategoryRelationRepository,
) SupplierFullService {
	return &supplierFullService{
		repoSupplier: repoSupplier,
		repoAddress:  repoAddress,
		repoContact:  repoContact,
		repoCatRel:   repoCatRel,
	}
}

func (s *supplierFullService) CreateFull(ctx context.Context, supplierFull *modelsFull.SupplierFull) (*modelsFull.SupplierFull, error) {
	if supplierFull == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := supplierFull.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

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

	createdSupplier, err := s.repoSupplier.CreateTx(ctx, tx, supplierFull.Supplier)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	supplierFull.Address.SupplierID = utils.StrToPtr(createdSupplier.ID)
	if err := supplierFull.Address.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("endereço inválido: %w", err))
	}
	createdAddress, err := s.repoAddress.CreateTx(ctx, tx, supplierFull.Address)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	supplierFull.Contact.SupplierID = utils.StrToPtr(createdSupplier.ID)
	if err := supplierFull.Contact.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("contato inválido: %w", err))
	}
	createdContact, err := s.repoContact.CreateTx(ctx, tx, supplierFull.Contact)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	for _, category := range supplierFull.Categories {
		relation := &modelsCatRel.SupplierCategoryRelations{
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

	result := &modelsFull.SupplierFull{
		Supplier:   createdSupplier,
		Address:    createdAddress,
		Contact:    createdContact,
		Categories: supplierFull.Categories,
	}

	return result, commitOrRollback(nil)
}
