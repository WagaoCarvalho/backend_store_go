package repositories

import (
	"context"

	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
)

type MockSupplierRepository struct {
	mock.Mock
}

func (m *MockSupplierRepository) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	args := m.Called(ctx, supplier)
	var result *models.Supplier
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Supplier)
	}
	return result, args.Error(1)
}

func (m *MockSupplierRepository) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	args := m.Called(ctx, id)
	var result *models.Supplier
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Supplier)
	}
	return result, args.Error(1)
}

func (m *MockSupplierRepository) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Supplier), args.Error(1)
}

func (m *MockSupplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	args := m.Called(ctx, supplier)
	return args.Error(0)
}

func (m *MockSupplierRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSupplierRepository) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSupplierRepository) Enable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSupplierRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}
