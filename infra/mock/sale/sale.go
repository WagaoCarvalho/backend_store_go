package mock

import (
	"context"
	"time"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockSale struct {
	mock.Mock
}

func (m *MockSale) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	args := m.Called(ctx, sale)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSale) CreateTx(ctx context.Context, tx pgx.Tx, sale *models.Sale) (*models.Sale, error) {
	args := m.Called(ctx, tx, sale)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSale) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSale) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, clientID, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSale) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, userID, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSale) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, status, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSale) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, start, end, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSale) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSale) Update(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSale) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSale) Cancel(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSale) Complete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSale) Activate(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSale) Returned(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSale) Filter(ctx context.Context, f *filter.SaleFilter) ([]*models.Sale, error) {
	args := m.Called(ctx, f)

	var result []*models.Sale
	if res := args.Get(0); res != nil {
		result = res.([]*models.Sale)
	}
	return result, args.Error(1)
}
