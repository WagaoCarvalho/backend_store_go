package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *supplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	suppliers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return suppliers, nil
}

func (s *supplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
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

func (s *supplierService) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
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

func (s *supplierService) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return version, nil
}

func (s *supplierService) SupplierExists(ctx context.Context, supplierID int64) (bool, error) {
	if supplierID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.repo.SupplierExists(ctx, supplierID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return exists, nil
}
