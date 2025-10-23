package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product_category"
	"github.com/stretchr/testify/mock"
)

type MockProductCategoryRepository struct {
	mock.Mock
}

func (m *MockProductCategoryRepository) GetAll(ctx context.Context) ([]*models.ProductCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepository) GetByID(ctx context.Context, id int64) (*models.ProductCategory, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepository) Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	args := m.Called(ctx, category)

	var result *models.ProductCategory
	if args.Get(0) != nil {
		result = args.Get(0).(*models.ProductCategory)
	}
	return result, args.Error(1)
}

func (m *MockProductCategoryRepository) Update(ctx context.Context, category *models.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
