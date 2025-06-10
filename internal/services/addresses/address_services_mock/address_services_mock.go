package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
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

func (m *MockAddressService) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockAddressService) Update(ctx context.Context, address *models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
