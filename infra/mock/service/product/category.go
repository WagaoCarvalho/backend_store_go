package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	"github.com/stretchr/testify/mock"
)

type ProductCategoryServiceMock struct {
	mock.Mock
}

func (m *ProductCategoryServiceMock) GetAll(ctx context.Context) ([]*models.ProductCategory, error) {
	args := m.Called(ctx)
	if categories, ok := args.Get(0).([]*models.ProductCategory); ok {
		return categories, args.Error(1)
	}
	return []*models.ProductCategory{}, args.Error(1)
}

func (m *ProductCategoryServiceMock) GetByID(ctx context.Context, id int64) (*models.ProductCategory, error) {
	args := m.Called(ctx, id)
	if category, ok := args.Get(0).(*models.ProductCategory); ok {
		return category, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductCategoryServiceMock) Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	args := m.Called(ctx, category)
	if created, ok := args.Get(0).(*models.ProductCategory); ok {
		return created, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductCategoryServiceMock) Update(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	args := m.Called(ctx, category)
	if updated, ok := args.Get(0).(*models.ProductCategory); ok {
		return updated, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductCategoryServiceMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
