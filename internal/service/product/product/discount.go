package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productService) DisableDiscount(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.DisableDiscount(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrProductDisableDiscount, err)
	}

	return nil
}

func (s *productService) EnableDiscount(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.EnableDiscount(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrProductEnableDiscount, err)
	}

	return nil
}

func (s *productService) ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if percent <= 0 {
		return nil, errMsg.ErrPercentInvalid
	}

	product, err := s.repo.ApplyDiscount(ctx, id, percent)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrProductApplyDiscount, err)
	}

	return product, nil
}
