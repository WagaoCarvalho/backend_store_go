package mock

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
	"github.com/stretchr/testify/mock"
)

type MockSupplierCategory struct {
	mock.Mock
}

func (m *MockSupplierCategory) Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error) {
	args := m.Called(ctx, category)
	result, _ := args.Get(0).(*models.SupplierCategory)
	return result, args.Error(1)
}

func (m *MockSupplierCategory) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategory) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategory) Update(ctx context.Context, category *models.SupplierCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockSupplierCategory) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
