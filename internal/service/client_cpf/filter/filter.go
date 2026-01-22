package services

import (
	"context"
	"fmt"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *clientCpfFiltertService) Filter(ctx context.Context, filter *filter.ClientCpfFilter) ([]*model.ClientCpf, error) {
	if filter == nil {
		return nil, fmt.Errorf("%w: filtro não pode ser nulo", errMsg.ErrInvalidFilter)
	}

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidFilter, err)
	}

	clients, err := s.repo.Filter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return clients, nil
}
