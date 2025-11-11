package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *saleItemService) Create(ctx context.Context, item *models.SaleItem) (*models.SaleItem, error) {
	if item == nil {
		return nil, errMsg.ErrInvalidData
	}

	if err := item.ValidateStructural(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := item.ValidateBusinessRules(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	createdItem, err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return createdItem, nil
}

func (s *saleItemService) Update(ctx context.Context, item *models.SaleItem) error {
	if item == nil {
		return errMsg.ErrInvalidData
	}

	if err := item.ValidateStructural(); err != nil {
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := item.ValidateBusinessRules(); err != nil {
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	return s.repo.Update(ctx, item)
}

func (s *saleItemService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	return s.repo.Delete(ctx, id)
}

func (s *saleItemService) DeleteBySaleID(ctx context.Context, saleID int64) error {
	if saleID <= 0 {
		return errMsg.ErrZeroID
	}

	return s.repo.DeleteBySaleID(ctx, saleID)
}
