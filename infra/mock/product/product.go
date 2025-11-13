package mock

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	"github.com/stretchr/testify/mock"
)

type ProductMock struct {
	mock.Mock
}

func (m *ProductMock) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductMock) GetAll(ctx context.Context) ([]*models.Product, error) {
	args := m.Called(ctx)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductMock) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	args := m.Called(ctx, id)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductMock) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductMock) GetVersionByID(ctx context.Context, pid int64) (int64, error) {
	args := m.Called(ctx, pid)

	if version, ok := args.Get(0).(int64); ok {
		return version, args.Error(1)
	}

	return 0, args.Error(1)
}

func (m *ProductMock) DisableProduct(ctx context.Context, pid int64) error {
	args := m.Called(ctx, pid)
	return args.Error(0)
}

func (m *ProductMock) EnableProduct(ctx context.Context, pid int64) error {
	args := m.Called(ctx, pid)
	return args.Error(0)
}

func (m *ProductMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductMock) UpdateStock(ctx context.Context, id int64, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *ProductMock) IncreaseStock(ctx context.Context, id int64, amount int) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *ProductMock) DecreaseStock(ctx context.Context, id int64, amount int) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *ProductMock) GetStock(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *ProductMock) EnableDiscount(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductMock) DisableDiscount(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductMock) ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error) {
	args := m.Called(ctx, id, percent)
	if product, ok := args.Get(0).(*models.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductMock) ProductExists(ctx context.Context, productID int64) (bool, error) {
	args := m.Called(ctx, productID)
	return args.Bool(0), args.Error(1)
}

func (m *ProductMock) Filter(ctx context.Context, filterData *models.ProductFilter) ([]*models.Product, error) {
	args := m.Called(ctx, filterData)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}
