package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *sale) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	if sale == nil {
		return nil, errMsg.ErrInvalidData
	}

	if err := sale.ValidateStructural(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := sale.ValidateBusinessRules(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	createdSale, err := s.repo.Create(ctx, sale)
	if err != nil {
		return nil, err
	}

	return createdSale, nil
}

func (s *sale) Update(ctx context.Context, sale *models.Sale) error {
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
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := sale.ValidateBusinessRules(); err != nil {
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := s.repo.Update(ctx, sale); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *sale) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
