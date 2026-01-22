package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *clientCPfService) Create(ctx context.Context, client *models.ClientCpf) (*models.ClientCpf, error) {
	if client == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	created, err := s.repo.Create(ctx, client)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *clientCPfService) Update(ctx context.Context, client *models.ClientCpf) error {
	if client.ID <= 0 {
		return errMsg.ErrZeroID
	}

	if client.Version <= 0 {
		return errMsg.ErrVersionConflict
	}

	if err := client.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if err := s.repo.Update(ctx, client); err != nil {
		if errors.Is(err, errMsg.ErrVersionConflict) ||
			errors.Is(err, errMsg.ErrDuplicate) ||
			errors.Is(err, errMsg.ErrInvalidData) {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *clientCPfService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}
	return nil
}
