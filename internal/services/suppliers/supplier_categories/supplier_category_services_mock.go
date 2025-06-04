package services

import (
	"context"

	supplier_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	"github.com/stretchr/testify/mock"
)

type MockSupplierCategoryService struct {
	mock.Mock
}

func (m *MockSupplierCategoryService) Create(ctx context.Context, category *supplier_categories.SupplierCategory) (int64, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierCategoryService) GetByID(ctx context.Context, id int64) (*supplier_categories.SupplierCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*supplier_categories.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryService) GetAll(ctx context.Context) ([]*supplier_categories.SupplierCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*supplier_categories.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryService) Update(ctx context.Context, category *supplier_categories.SupplierCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockSupplierCategoryService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
