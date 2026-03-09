package services

import (
	"context"
	"fmt"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/address/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *addressFiltertService) Filter(ctx context.Context, filter *filter.AddressFilter) ([]*model.Address, error) {
	if filter == nil {
		return nil, fmt.Errorf("%w: filtro não pode ser nulo", errMsg.ErrInvalidFilter)
	}

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidFilter, err)
	}

	address, err := s.repo.Filter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return address, nil
}
