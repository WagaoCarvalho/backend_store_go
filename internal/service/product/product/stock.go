package services

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productService) UpdateStock(ctx context.Context, id int64, quantity int) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if quantity <= 0 {
		return errMsg.ErrInvalidQuantity
	}

	err := s.repo.UpdateStock(ctx, id, quantity)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *productService) IncreaseStock(ctx context.Context, id int64, amount int) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if amount <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.IncreaseStock(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil

}

func (s *productService) DecreaseStock(ctx context.Context, id int64, amount int) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if amount <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.DecreaseStock(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *productService) GetStock(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return 0, errMsg.ErrZeroID
	}

	stock, err := s.repo.GetStock(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return stock, nil
}
