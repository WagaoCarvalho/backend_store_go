package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
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
	repo repo.SupplierRepository
}

func NewSupplierService(repo repo.SupplierRepository) SupplierService {
	return &supplierService{
		repo: repo,
	}
}

func (s *supplierService) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	if supplier == nil || supplier.Name == "" {
		return nil, err_msg.ErrInvalidData
	}

	created, err := s.repo.Create(ctx, supplier)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
	}

	return created, nil
}

func (s *supplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	suppliers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return suppliers, nil
}

func (s *supplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	if id <= 0 {
		return nil, err_msg.ErrID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	if supplier == nil {
		return nil, err_msg.ErrNotFound
	}

	return supplier, nil
}

func (s *supplierService) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	if name == "" {
		return nil, err_msg.ErrInvalidData
	}

	suppliers, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	if len(suppliers) == 0 {
		return nil, err_msg.ErrNotFound
	}

	return suppliers, nil
}

func (s *supplierService) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, err_msg.ErrID
	}

	version, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return version, nil
}

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	if supplier.ID <= 0 {
		return nil, err_msg.ErrID
	}

	if supplier.Name == "" {
		return nil, err_msg.ErrInvalidData
	}

	if supplier.Version == 0 {
		return nil, err_msg.ErrVersionConflict
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		switch {
		case errors.Is(err, err_msg.ErrVersionConflict):
			return nil, err_msg.ErrVersionConflict
		case errors.Is(err, err_msg.ErrNotFound):
			return nil, err_msg.ErrNotFound
		default:
			return nil, fmt.Errorf("%w: %v", err_msg.ErrUpdate, err)
		}
	}

	return supplier, nil
}

func (s *supplierService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return err_msg.ErrID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *supplierService) Disable(ctx context.Context, id int64) error {
	if id <= 0 {
		return err_msg.ErrID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	supplier.Status = false

	if err := s.repo.Update(ctx, supplier); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDisable, err)
	}

	return nil
}

func (s *supplierService) Enable(ctx context.Context, id int64) error {
	if id <= 0 {
		return err_msg.ErrID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	supplier.Status = true

	if err := s.repo.Update(ctx, supplier); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrEnable, err)
	}

	return nil
}
