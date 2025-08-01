package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/suppliers"
)

type SupplierService interface {
	Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetByName(ctx context.Context, name string) ([]*models.Supplier, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	Update(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error)
	Delete(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
}

type supplierService struct {
	repo   repository.SupplierRepository
	logger *logger.LoggerAdapter
}

func NewSupplierService(repo repository.SupplierRepository, logger *logger.LoggerAdapter) SupplierService {
	return &supplierService{
		repo:   repo,
		logger: logger,
	}
}

func (s *supplierService) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	ref := "[supplierService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name":   supplier.Name,
		"cnpj":   supplier.CNPJ,
		"status": supplier.Status,
	})

	if supplier == nil || supplier.Name == "" {
		s.logger.Error(ctx, ErrSupplierNameRequired, ref+"nome do fornecedor é obrigatório", nil)
		return nil, ErrSupplierNameRequired
	}

	created, err := s.repo.Create(ctx, supplier)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": supplier.Name,
			"cnpj": supplier.CNPJ,
		})
		return nil, fmt.Errorf("%w: %v", ErrSupplierCreateFailed, err)
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": created.ID,
		"name":        created.Name,
		"cnpj":        created.CNPJ,
	})

	return created, nil
}

func (s *supplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	ref := "[supplierService - GetAll] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, nil)

	suppliers, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrGetSupplier, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(suppliers),
	})
	return suppliers, nil
}

func (s *supplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	ref := "[supplierService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": id,
	})

	if id <= 0 {
		s.logger.Error(ctx, ErrInvalidSupplierID, ref+"ID inválido", map[string]any{
			"supplier_id": id,
		})
		return nil, ErrInvalidSupplierID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetSupplier, err)
	}

	if supplier == nil {
		s.logger.Warn(ctx, ref+"fornecedor não encontrado", map[string]any{
			"supplier_id": id,
		})
		return nil, ErrSupplierNotFound
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplier.ID,
		"name":        supplier.Name,
	})
	return supplier, nil
}

func (s *supplierService) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	ref := "[supplierService - GetByName] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_name": name,
	})

	if name == "" {
		s.logger.Error(ctx, ErrInvalidSupplierName, ref+"Nome inválido", map[string]any{
			"supplier_name": name,
		})
		return nil, ErrInvalidSupplierName
	}

	suppliers, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_name": name,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetSupplier, err)
	}

	if len(suppliers) == 0 {
		s.logger.Warn(ctx, ref+"fornecedor não encontrado", map[string]any{
			"supplier_name": name,
		})
		return nil, ErrSupplierNotFound
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(suppliers),
	})

	return suppliers, nil
}

func (s *supplierService) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	ref := "[supplierService - GetVersionByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": id,
	})

	if id <= 0 {
		s.logger.Error(ctx, ErrInvalidSupplierID, ref+"ID inválido", map[string]any{
			"supplier_id": id,
		})
		return 0, ErrInvalidSupplierID
	}

	version, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		return 0, fmt.Errorf("%w: %v", ErrGetSupplierVersion, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": id,
		"version":     version,
	})
	return version, nil
}

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	ref := "[supplierService - Update] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"supplier_id": supplier.ID,
	})

	if supplier.ID <= 0 {
		s.logger.Error(ctx, ErrInvalidSupplierID, ref+"ID inválido", map[string]any{
			"supplier_id": supplier.ID,
		})
		return nil, ErrInvalidSupplierID
	}

	if supplier.Name == "" {
		s.logger.Error(ctx, ErrSupplierNameRequired, ref+"nome do fornecedor obrigatório", nil)
		return nil, ErrSupplierNameRequired
	}

	if supplier.Version == 0 {
		s.logger.Error(ctx, ErrSupplierVersionRequired, ref+"versão do fornecedor obrigatória", nil)
		return nil, ErrSupplierVersionRequired
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		switch {
		case errors.Is(err, repository.ErrVersionConflict):
			s.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"supplier_id": supplier.ID,
				"version":     supplier.Version,
			})
			return nil, ErrSupplierVersionConflict

		case errors.Is(err, repository.ErrSupplierNotFound):
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": supplier.ID,
			})
			return nil, ErrSupplierNotFound

		default:
			s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"supplier_id": supplier.ID,
			})
			return nil, fmt.Errorf("%w: %v", ErrSupplierUpdate, err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": supplier.ID,
	})

	return supplier, nil
}

func (s *supplierService) Delete(ctx context.Context, id int64) error {
	ref := "[supplierService - Delete] - "
	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"supplier_id": id,
	})

	if id <= 0 {
		s.logger.Error(ctx, ErrInvalidSupplierIDForDeletion, ref+"ID inválido para deleção", map[string]any{
			"supplier_id": id,
		})
		return ErrInvalidSupplierIDForDeletion
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"supplier_id": id,
		})
		return err
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": id,
	})

	return nil
}

func (s *supplierService) Disable(ctx context.Context, id int64) error {
	ref := "[supplierService - Disable] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"supplier_id": id,
	})

	if id <= 0 {
		s.logger.Error(ctx, ErrInvalidSupplierID, ref+"ID inválido para desabilitar fornecedor", map[string]any{
			"supplier_id": id,
		})
		return ErrInvalidSupplierID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		return fmt.Errorf("%w: %v", ErrGetSupplier, err)
	}

	supplier.Status = false

	err = s.repo.Update(ctx, supplier)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"supplier_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDisableSupplier, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
		"status":      supplier.Status,
	})

	return nil
}

func (s *supplierService) Enable(ctx context.Context, id int64) error {
	ref := "[supplierService - Enable] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"supplier_id": id,
	})

	if id <= 0 {
		s.logger.Error(ctx, ErrInvalidSupplierID, ref+"ID inválido para habilitar fornecedor", map[string]any{
			"supplier_id": id,
		})
		return ErrInvalidSupplierID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		return fmt.Errorf("%w: %v", ErrGetSupplier, err)
	}

	supplier.Status = true

	err = s.repo.Update(ctx, supplier)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"supplier_id": id,
		})
		return fmt.Errorf("%w: %v", ErrEnableSupplier, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
		"status":      supplier.Status,
	})

	return nil
}
