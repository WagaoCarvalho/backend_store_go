package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	"github.com/stretchr/testify/mock"
)

type MockSupplierService struct {
	mock.Mock
}

func (m *MockSupplierService) Create(ctx context.Context, s *models.Supplier) (*models.Supplier, error) {
	args := m.Called(ctx, s)
	if supplier, ok := args.Get(0).(*models.Supplier); ok {
		return supplier, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSupplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	args := m.Called(ctx)
	if suppliers, ok := args.Get(0).([]*models.Supplier); ok {
		return suppliers, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSupplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	args := m.Called(ctx, id)
	if supplier, ok := args.Get(0).(*models.Supplier); ok {
		return supplier, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSupplierService) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Supplier), args.Error(1)
}

func (m *MockSupplierService) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierService) Update(ctx context.Context, s *models.Supplier) (*models.Supplier, error) {
	args := m.Called(ctx, s)
	if supplier, ok := args.Get(0).(*models.Supplier); ok {
		return supplier, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSupplierService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSupplierService) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSupplierService) Enable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
