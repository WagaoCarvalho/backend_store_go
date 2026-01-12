package services

import (
	"context"
	"fmt"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *supplierFilterService) Filter(
	ctx context.Context,
	filter *filter.SupplierFilter,
) ([]*model.Supplier, error) {

	if filter == nil {
		return nil, fmt.Errorf("%w: filtro não pode ser nulo", errMsg.ErrInvalidFilter)
	}

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidFilter, err)
	}

	suppliers, err := s.repo.Filter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return suppliers, nil
}
