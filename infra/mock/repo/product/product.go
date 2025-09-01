package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	args := m.Called(ctx, limit, offset)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	args := m.Called(ctx, id)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
	args := m.Called(ctx, name)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {
	args := m.Called(ctx, manufacturer)
	if prods, ok := args.Get(0).([]*models.Product); ok {
		return prods, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) GetVersionByID(ctx context.Context, pid int64) (int64, error) {
	args := m.Called(ctx, pid)

	if version, ok := args.Get(0).(int64); ok {
		return version, args.Error(1)
	}

	return 0, args.Error(1)
}

func (m *ProductRepositoryMock) DisableProduct(ctx context.Context, pid int64) error {
	args := m.Called(ctx, pid)
	return args.Error(0)
}

func (m *ProductRepositoryMock) EnableProduct(ctx context.Context, pid int64) error {
	args := m.Called(ctx, pid)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) UpdateStock(ctx context.Context, id int64, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *ProductRepositoryMock) IncreaseStock(ctx context.Context, id int64, amount int) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *ProductRepositoryMock) DecreaseStock(ctx context.Context, id int64, amount int) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetStock(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *ProductRepositoryMock) EnableDiscount(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) DisableDiscount(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error) {
	args := m.Called(ctx, id, percent)
	if product, ok := args.Get(0).(*models.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}
