package supplier_service_mocks

import (
	"context"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	"github.com/stretchr/testify/mock"
)

type MockAddressService struct {
	mock.Mock
}

func (m *MockAddressService) Create(ctx context.Context, address *models_address.Address) (*models_address.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*models_address.Address), args.Error(1)
}

func (m *MockAddressService) GetByID(ctx context.Context, id int64) (*models_address.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models_address.Address), args.Error(1)
}

func (m *MockAddressService) Update(ctx context.Context, address *models_address.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressService) Delete(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}
