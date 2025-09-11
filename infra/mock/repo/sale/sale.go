package repositories

import (
	"context"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockSaleRepository struct {
	mock.Mock
}

func (m *MockSaleRepository) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	args := m.Called(ctx, sale)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) CreateTx(ctx context.Context, tx pgx.Tx, sale *models.Sale) (*models.Sale, error) {
	args := m.Called(ctx, tx, sale)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, clientID, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, userID, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, status, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, start, end, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) Update(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
