package services

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *clientService) Disable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Disable(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}

func (s *clientService) Enable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Enable(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}
