package services

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productService) DisableDiscount(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.DisableDiscount(ctx, id); err != nil {
		if err == errMsg.ErrNotFound {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductDisableDiscount, err)
	}

	return nil
}

func (s *productService) EnableDiscount(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.EnableDiscount(ctx, id); err != nil {
		if err == errMsg.ErrNotFound {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductEnableDiscount, err)
	}

	return nil
}

func (s *productService) ApplyDiscount(ctx context.Context, id int64, percent float64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	// Validação completa do percentual
	if percent < 0 || percent > 100 {
		return errMsg.ErrInvalidDiscountPercent
	}

	if err := s.repo.ApplyDiscount(ctx, id, percent); err != nil {
		// Propaga erros específicos sem alterá-los
		if err == errMsg.ErrNotFound ||
			err == errMsg.ErrProductDiscountNotAllowed ||
			err == errMsg.ErrInvalidDiscountPercent {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductApplyDiscount, err)
	}

	return nil
}
