package services

import (
	"context"
	"fmt"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userFilterService) Filter(
	ctx context.Context,
	filter *filter.UserFilter,
) ([]*model.User, error) {

	if filter == nil {
		return nil, fmt.Errorf("%w: filtro não pode ser nulo", errMsg.ErrInvalidFilter)
	}

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidFilter, err)
	}

	users, err := s.repo.Filter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return users, nil
}
