package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
)

type MockSaleItem struct {
	mock.Mock
}

func (m *MockSaleItem) GetByID(ctx context.Context, id int64) (*models.SaleItem, error) {
	args := m.Called(ctx, id)
	if item, ok := args.Get(0).(*models.SaleItem); ok {
		return item, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleItem) GetBySaleID(ctx context.Context, saleID int64, limit, offset int) ([]*models.SaleItem, error) {
	args := m.Called(ctx, saleID, limit, offset)
	if items, ok := args.Get(0).([]*models.SaleItem); ok {
		return items, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleItem) GetByProductID(ctx context.Context, productID int64, limit, offset int) ([]*models.SaleItem, error) {
	args := m.Called(ctx, productID, limit, offset)
	if items, ok := args.Get(0).([]*models.SaleItem); ok {
		return items, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleItem) Create(ctx context.Context, item *models.SaleItem) (*models.SaleItem, error) {
	args := m.Called(ctx, item)
	if created, ok := args.Get(0).(*models.SaleItem); ok {
		return created, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleItem) Update(ctx context.Context, item *models.SaleItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockSaleItem) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSaleItem) DeleteBySaleID(ctx context.Context, saleID int64) error {
	args := m.Called(ctx, saleID)
	return args.Error(0)
}

func (m *MockSaleItem) ItemExists(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
