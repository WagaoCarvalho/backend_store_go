package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockAddressRepository struct {
	mock.Mock
}

func (m *MockAddressRepository) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressRepository) CreateTx(ctx context.Context, tx pgx.Tx, address *models.Address) (*models.Address, error) {
	args := m.Called(ctx, tx, address)

	// Protege contra retorno nil em testes
	var result *models.Address
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Address)
	}

	return result, args.Error(1)
}

func (m *MockAddressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetByUserID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).([]*models.Address); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressRepository) GetByClientID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).([]*models.Address); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressRepository) GetBySupplierID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).([]*models.Address); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressRepository) Update(ctx context.Context, address *models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
