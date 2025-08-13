package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/stretchr/testify/mock"
)

type MockAddressService struct {
	mock.Mock
}

func (m *MockAddressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressService) GetByUserID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)

	// Precaução para evitar panic se o valor retornado for nil
	if addresses, ok := args.Get(0).([]*models.Address); ok {
		return addresses, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressService) GetByClientID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)

	if addresses, ok := args.Get(0).([]*models.Address); ok {
		return addresses, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressService) GetBySupplierID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)

	if addresses, ok := args.Get(0).([]*models.Address); ok {
		return addresses, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressService) Update(ctx context.Context, address *models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
