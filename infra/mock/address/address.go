package mock

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockAddress struct {
	mock.Mock
}

func (m *MockAddress) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	args := m.Called(ctx, address)
	if result := args.Get(0); result != nil {
		return result.(*models.Address), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddress) CreateTx(ctx context.Context, tx pgx.Tx, address *models.Address) (*models.Address, error) {
	args := m.Called(ctx, tx, address)

	// Protege contra retorno nil em testes
	var result *models.Address
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Address)
	}

	return result, args.Error(1)
}

func (m *MockAddress) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (m *MockAddress) GetByUserID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).([]*models.Address); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddress) GetByClientCpfID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).([]*models.Address); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddress) GetBySupplierID(ctx context.Context, id int64) ([]*models.Address, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).([]*models.Address); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddress) Update(ctx context.Context, address *models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddress) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAddress) Disable(ctx context.Context, aid int64) error {
	args := m.Called(ctx, aid)
	return args.Error(0)
}

func (m *MockAddress) Enable(ctx context.Context, aid int64) error {
	args := m.Called(ctx, aid)
	return args.Error(0)
}
