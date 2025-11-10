package services

import (
	"context"
	"fmt"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *client) GetAll(ctx context.Context, filter *model.ClientFilter) ([]*model.Client, error) {
	if filter == nil {
		return nil, fmt.Errorf("%w: filtro não pode ser nulo", errMsg.ErrInvalidFilter)
	}

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidFilter, err)
	}

	clients, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return clients, nil
}
