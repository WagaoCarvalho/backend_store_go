package services

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *product) DisableProduct(ctx context.Context, uid int64) error {

	if uid <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.DisableProduct(ctx, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (s *product) EnableProduct(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.EnableProduct(ctx, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}
