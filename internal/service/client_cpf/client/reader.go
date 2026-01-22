package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *clientCPfService) GetByID(ctx context.Context, id int64) (*models.ClientCpf, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return client, nil
}
