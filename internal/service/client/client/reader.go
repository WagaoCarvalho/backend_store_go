package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *clientService) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return client, nil
}

func (s *clientService) GetByName(ctx context.Context, name string) ([]*models.Client, error) {
	if name == "" {
		return nil, errMsg.ErrInvalidData
	}

	clients, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return clients, nil
}
