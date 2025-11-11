package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *supplierService) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	if err := supplier.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	created, err := s.repo.Create(ctx, supplier)
	if err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrCreate)
	}

	return created, nil
}

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) error {
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

func (s *supplierService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
