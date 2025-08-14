package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/supplier/supplier_categories"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type SupplierCategoryService interface {
	Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error)
	GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error)
	GetAll(ctx context.Context) ([]*models.SupplierCategory, error)
	Update(ctx context.Context, category *models.SupplierCategory) error
	Delete(ctx context.Context, id int64) error
}

type supplierCategoryService struct {
	repo   repository.SupplierCategoryRepository
	logger *logger.LoggerAdapter
}

func NewSupplierCategoryService(repo repository.SupplierCategoryRepository, logger *logger.LoggerAdapter) SupplierCategoryService {
	return &supplierCategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *supplierCategoryService) Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error) {
	ref := "[supplierCategoryService - Create] - "

	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name": category.Name,
	})

	if err := category.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"name":  category.Name,
			"error": err.Error(),
		})
		return nil, err
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": category.Name,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"category_id": createdCategory.ID,
		"name":        createdCategory.Name,
	})

	return createdCategory, nil
}

func (s *supplierCategoryService) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	ref := "[supplierCategoryService - GetByID] - "

	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"id": id,
		})
		return nil, ErrCategoryIDInvalid
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetCategory, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"id": category.ID,
	})

	return category, nil
}

func (s *supplierCategoryService) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	ref := "[supplierCategoryService - GetAll] - "

	s.logger.Info(ctx, ref+logger.LogGetInit, nil)

	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrGetCategories, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(categories),
	})

	return categories, nil
}

func (s *supplierCategoryService) Update(ctx context.Context, category *models.SupplierCategory) error {
	ref := "[supplierCategoryService - Update] - "

	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"id":   category.ID,
		"name": category.Name,
	})

	if category.ID == 0 {
		s.logger.Warn(ctx, ref+"ValidationError - id zero", nil)
		return ErrCategoryIDRequired
	}

	if err := category.Validate(); err != nil {
		s.logger.Warn(ctx, ref+"ValidationError", map[string]any{
			"id":    category.ID,
			"error": err.Error(),
		})
		return err
	}

	if err := s.repo.Update(ctx, category); err != nil {
		s.logger.Error(ctx, err, ref+"UpdateError", map[string]any{
			"id": category.ID,
		})
		return fmt.Errorf("%w: %v", ErrUpdateCategory, err)
	}

	s.logger.Info(ctx, ref+"UpdateSuccess", map[string]any{
		"id": category.ID,
	})

	return nil
}

func (s *supplierCategoryService) Delete(ctx context.Context, id int64) error {
	ref := "[supplierCategoryService - Delete] - "

	s.logger.Info(ctx, ref+"Init", map[string]any{
		"id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+"ValidationError - invalid id", map[string]any{
			"id": id,
		})
		return ErrCategoryDeleteInvalidID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, err, ref+"DeleteError", map[string]any{
			"id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}

	s.logger.Info(ctx, ref+"DeleteSuccess", map[string]any{
		"id": id,
	})

	return nil
}
