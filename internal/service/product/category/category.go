package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category"
)

type ProductCategoryService interface {
	GetAll(ctx context.Context) ([]*models.ProductCategory, error)
	GetByID(ctx context.Context, id int64) (*models.ProductCategory, error)
	Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error)
	Update(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error)
	Delete(ctx context.Context, id int64) error
}

type productCategoryService struct {
	repo repo.ProductCategoryRepository
}

func NewProductCategoryService(repo repo.ProductCategoryRepository) ProductCategoryService {
	return &productCategoryService{
		repo: repo,
	}
}

func (s *productCategoryService) Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return createdCategory, nil
}

func (s *productCategoryService) GetAll(ctx context.Context) ([]*models.ProductCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return categories, nil
}

func (s *productCategoryService) GetByID(ctx context.Context, id int64) (*models.ProductCategory, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return category, nil
}

func (s *productCategoryService) Update(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	if category.ID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.repo.GetByID(ctx, int64(category.ID)); err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return category, nil
}

func (s *productCategoryService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
