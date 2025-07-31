package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockSupplierCategoryRepo struct {
	mock.Mock
}

func (m *MockSupplierCategoryRepo) Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error) {
	args := m.Called(ctx, category)
	result, _ := args.Get(0).(*models.SupplierCategory)
	return result, args.Error(1)
}

func (m *MockSupplierCategoryRepo) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryRepo) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryRepo) Update(ctx context.Context, category *models.SupplierCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockSupplierCategoryRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
