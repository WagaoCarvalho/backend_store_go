package services

import (
	"context"
	"errors"
	"fmt"

	models_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	models_full "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	repo_address "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repo_contact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repo_relation "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category_relations"
	repo_supplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_full_repositories"
)

type SupplierFullService interface {
	CreateFull(ctx context.Context, supplierFull *models_full.SupplierFull) (*models_full.SupplierFull, error)
}

type supplierFullService struct {
	repo_supplier repo_supplier.SupplierFullRepository
	repo_address  repo_address.AddressRepository
	repo_contact  repo_contact.ContactRepository
	repo_cat_rel  repo_relation.SupplierCategoryRelationRepository
	logger        *logger.LoggerAdapter
}

func NewSupplierFullService(
	repo_supplier repo_supplier.SupplierFullRepository,
	repo_address repo_address.AddressRepository,
	repo_contact repo_contact.ContactRepository,
	repo_cat_rel repo_relation.SupplierCategoryRelationRepository,
	logger *logger.LoggerAdapter,
) SupplierFullService {
	return &supplierFullService{
		repo_supplier: repo_supplier,
		repo_address:  repo_address,
		repo_contact:  repo_contact,
		repo_cat_rel:  repo_cat_rel,
		logger:        logger,
	}
}

func (s *supplierFullService) CreateFull(ctx context.Context, supplierFull *models_full.SupplierFull) (*models_full.SupplierFull, error) {
	const ref = "[supplierService - CreateFull] - "

	logFields := map[string]any{}
	if supplierFull != nil && supplierFull.Supplier != nil {
		logFields["name"] = supplierFull.Supplier.Name
		logFields["cpf"] = supplierFull.Supplier.CPF
		logFields["cnpj"] = supplierFull.Supplier.CNPJ
	}

	s.logger.Info(ctx, ref+logger.LogCreateInit, logFields)

	if err := supplierFull.Validate(); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogValidateError, logFields)
		return nil, err
	}

	tx, err := s.repo_supplier.BeginTx(ctx)
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

	createdSupplier, err := s.repo_supplier.CreateTx(ctx, tx, supplierFull.Supplier)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	supplierFull.Address.SupplierID = utils.StrToPtr(createdSupplier.ID)

	if err := supplierFull.Address.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("endereço inválido: %w", err))
	}

	createdAddress, err := s.repo_address.CreateTx(ctx, tx, supplierFull.Address)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	supplierFull.Contact.SupplierID = utils.StrToPtr(createdSupplier.ID)

	if err := supplierFull.Contact.Validate(); err != nil {
		return nil, commitOrRollback(fmt.Errorf("contato inválido: %w", err))
	}

	createdContact, err := s.repo_contact.CreateTx(ctx, tx, supplierFull.Contact)
	if err != nil {
		return nil, commitOrRollback(err)
	}

	for _, category := range supplierFull.Categories {
		relation := &models_cat_rel.SupplierCategoryRelations{
			SupplierID: createdSupplier.ID,
			CategoryID: int64(category.ID),
		}

		if err := relation.Validate(); err != nil {
			return nil, commitOrRollback(fmt.Errorf("relação fornecedor-categoria inválida: %w", err))
		}

		if _, err := s.repo_cat_rel.CreateTx(ctx, tx, relation); err != nil {
			return nil, commitOrRollback(err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdSupplier.ID,
		"name":        createdSupplier.Name,
	})

	result := &models_full.SupplierFull{
		Supplier:   createdSupplier,
		Address:    createdAddress,
		Contact:    createdContact,
		Categories: supplierFull.Categories,
	}

	return result, commitOrRollback(nil)
}
