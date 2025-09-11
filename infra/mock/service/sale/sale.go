package services

import (
	"context"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	"github.com/stretchr/testify/mock"
)

type MockSaleService struct {
	mock.Mock
}

func (m *MockSaleService) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	args := m.Called(ctx, sale)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleService) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleService) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, clientID, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleService) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, userID, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleService) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, status, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleService) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	args := m.Called(ctx, start, end, limit, offset, orderBy, orderDir)
	if result := args.Get(0); result != nil {
		return result.([]*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleService) Update(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
