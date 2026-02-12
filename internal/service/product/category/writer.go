package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productCategoryService) Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	if err := category.Validate(); err != nil {
		return nil, err // Retorna o erro específico de validação
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		if errors.Is(err, errMsg.ErrDuplicate) { // Consistência: ErrDuplicate
			return nil, fmt.Errorf("%w: %v", errMsg.ErrDuplicate, err)
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return createdCategory, nil
}

func (s *productCategoryService) Update(ctx context.Context, category *models.ProductCategory) error {
	if category.ID <= 0 {
		return errMsg.ErrZeroID
	}

	if err := category.Validate(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, category); err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *productCategoryService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	// Remove a verificação duplicada - o repo.Delete já verifica
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
