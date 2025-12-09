package services

import (
	"context"
	"fmt"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productService) Filter(ctx context.Context, filterData *modelFilter.ProductFilter) ([]*models.Product, error) {
	if filterData == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := filterData.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	products, err := s.repo.Filter(ctx, filterData)
	if err != nil {
		return nil, err
	}

	return products, nil
}
