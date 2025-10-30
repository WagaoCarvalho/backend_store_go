package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *product) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	if product == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := product.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (s *product) Update(ctx context.Context, product *models.Product) error {

	if product.ID <= 0 {
		return errMsg.ErrZeroID
	}
	if err := product.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if product.Version <= 0 {
		return errMsg.ErrVersionConflict
	}

	err := s.repo.Update(ctx, product)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			exists, errCheck := s.repo.ProductExists(ctx, product.ID)
			if errCheck != nil {
				return fmt.Errorf("%w: %v", errMsg.ErrGet, errCheck)
			}

			if !exists {
				return errMsg.ErrNotFound
			}
			return errMsg.ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *product) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
