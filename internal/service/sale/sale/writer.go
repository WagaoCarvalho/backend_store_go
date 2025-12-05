package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *saleService) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	if sale == nil {
		return nil, errMsg.ErrInvalidData
	}

	if err := sale.ValidateStructural(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if err := sale.ValidateBusinessRules(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	createdSale, err := s.repo.Create(ctx, sale)
	if err != nil {
		return nil, err
	}

	return createdSale, nil
}

func (s *saleService) Update(ctx context.Context, sale *models.Sale) error {
	if sale == nil {
		return errMsg.ErrInvalidData
	}
	if sale.ID <= 0 {
		return errMsg.ErrZeroID
	}
	if sale.Version <= 0 {
		return errMsg.ErrVersionConflict
	}

	if err := sale.ValidateStructural(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if err := sale.ValidateBusinessRules(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	err := s.repo.Update(ctx, sale)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			return errMsg.ErrNotFound
		case errors.Is(err, errMsg.ErrVersionConflict):
			return errMsg.ErrVersionConflict
		default:
			return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (s *saleService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
