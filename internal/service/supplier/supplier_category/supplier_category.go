package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category"
)

type SupplierCategory interface {
	Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error)
	GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error)
	GetAll(ctx context.Context) ([]*models.SupplierCategory, error)
	Update(ctx context.Context, category *models.SupplierCategory) error
	Delete(ctx context.Context, id int64) error
}

type supplierCategory struct {
	repo repo.SupplierCategory
}

func NewSupplierCategory(repo repo.SupplierCategory) SupplierCategory {
	return &supplierCategory{
		repo: repo,
	}
}

func (s *supplierCategory) Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error) {
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return createdCategory, nil
}

func (s *supplierCategory) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return category, nil
}

func (s *supplierCategory) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return categories, nil
}

func (s *supplierCategory) Update(ctx context.Context, category *models.SupplierCategory) error {
	if category.ID <= 0 {
		return errMsg.ErrZeroID
	}

	if err := category.Validate(); err != nil {
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *supplierCategory) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
