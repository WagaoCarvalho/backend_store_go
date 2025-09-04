package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_categories"
)

type SupplierCategoryService interface {
	Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error)
	GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error)
	GetAll(ctx context.Context) ([]*models.SupplierCategory, error)
	Update(ctx context.Context, category *models.SupplierCategory) error
	Delete(ctx context.Context, id int64) error
}

type supplierCategoryService struct {
	repo repo.SupplierCategoryRepository
}

func NewSupplierCategoryService(repo repo.SupplierCategoryRepository) SupplierCategoryService {
	return &supplierCategoryService{
		repo: repo,
	}
}

func (s *supplierCategoryService) Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error) {
	if err := category.Validate(); err != nil {
		return nil, err
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
	}

	return createdCategory, nil
}

func (s *supplierCategoryService) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	if id <= 0 {
		return nil, err_msg.ErrID
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return category, nil
}

func (s *supplierCategoryService) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return categories, nil
}

func (s *supplierCategoryService) Update(ctx context.Context, category *models.SupplierCategory) error {
	if category.ID == 0 {
		return err_msg.ErrID
	}

	if err := category.Validate(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrUpdate, err)
	}

	return nil
}

func (s *supplierCategoryService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return err_msg.ErrID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
