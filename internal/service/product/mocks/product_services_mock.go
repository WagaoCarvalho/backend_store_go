package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	args := m.Called(ctx, limit, offset)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) GetById(ctx context.Context, id int64) (*models.Product, error) {
	args := m.Called(ctx, id)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
	args := m.Called(ctx, name)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {
	args := m.Called(ctx, manufacturer)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(int64), args.Error(1)
}

func (m *ProductServiceMock) DisableProduct(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *ProductServiceMock) EnableProduct(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *ProductServiceMock) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductServiceMock) UpdateStock(ctx context.Context, id int64, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *ProductServiceMock) IncreaseStock(ctx context.Context, id int64, amount int) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *ProductServiceMock) DecreaseStock(ctx context.Context, id int64, amount int) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *ProductServiceMock) GetStock(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *ProductServiceMock) EnableDiscount(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
