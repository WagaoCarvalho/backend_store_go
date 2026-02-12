package services

import (
	"context"
	"fmt"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

// No serviço:
func (s *productFilterService) Filter(ctx context.Context, filter *filter.ProductFilter) ([]*model.Product, error) {
	if filter == nil {
		return nil, errMsg.ErrInvalidFilter
	}

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidFilter, err)
	}

	products, err := s.repo.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return products, nil
}
