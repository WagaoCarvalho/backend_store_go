package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
)

type Supplier interface {
	Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetByName(ctx context.Context, name string) ([]*models.Supplier, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
}

type supplier struct {
	repo repo.Supplier
}

func NewSupplier(repo repo.Supplier) Supplier {
	return &supplier{
		repo: repo,
	}
}

func (s *supplier) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	if err := supplier.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	created, err := s.repo.Create(ctx, supplier)
	if err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrCreate)
	}

	return created, nil
}

func (s *supplier) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	suppliers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return suppliers, nil
}

func (s *supplier) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if supplier == nil {
		return nil, errMsg.ErrNotFound
	}

	return supplier, nil
}

func (s *supplier) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	if name == "" {
		return nil, errMsg.ErrInvalidData
	}

	suppliers, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(suppliers) == 0 {
		return nil, errMsg.ErrNotFound
	}

	return suppliers, nil
}

func (s *supplier) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return version, nil
}

func (s *supplier) Update(ctx context.Context, supplier *models.Supplier) error {
	if supplier.ID <= 0 {
		return errMsg.ErrZeroID
	}

	if err := supplier.Validate(); err != nil {
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if supplier.Version == 0 {
		return errMsg.ErrVersionConflict
	}

	err := s.repo.Update(ctx, supplier)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			exists, errCheck := s.repo.SupplierExists(ctx, supplier.ID)
			if errCheck != nil {
				return fmt.Errorf("%w: %v", errMsg.ErrGet, errCheck)
			}

			if !exists {
				return errMsg.ErrNotFound
			}
			return errMsg.ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *supplier) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *supplier) Disable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	supplier.Status = false

	if err := s.repo.Update(ctx, supplier); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (s *supplier) Enable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	supplier.Status = true

	if err := s.repo.Update(ctx, supplier); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}
