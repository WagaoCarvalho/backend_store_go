package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	"github.com/stretchr/testify/mock"
)

type MockAddressRepository struct {
	mock.Mock
}

func (m *MockAddressRepository) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetByUserID(ctx context.Context, id int64) (*models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetByClientID(ctx context.Context, id int64) (*models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetBySupplierID(ctx context.Context, id int64) (*models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}
func (m *MockAddressRepository) Update(ctx context.Context, address *models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
