package services

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productService) DisableProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.DisableProduct(ctx, id); err != nil {
		// Propaga erro NotFound sem alterar
		if err == errMsg.ErrNotFound {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (s *productService) EnableProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.EnableProduct(ctx, id); err != nil {
		// Propaga erro NotFound sem alterar
		if err == errMsg.ErrNotFound {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}
