package mock

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// Mock do repositório
type MockSupplierCategoryRelation struct {
	mock.Mock
}

func (m *MockSupplierCategoryRelation) Create(ctx context.Context, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
	args := m.Called(ctx, relation)
	if rel, ok := args.Get(0).(*models.SupplierCategoryRelation); ok {
		return rel, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSupplierCategoryRelation) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
	args := m.Called(ctx, tx, relation)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.SupplierCategoryRelation), args.Error(1)
}

func (m *MockSupplierCategoryRelation) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models.SupplierCategoryRelation), args.Error(1)
}

func (m *MockSupplierCategoryRelation) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*models.SupplierCategoryRelation), args.Error(1)
}

func (m *MockSupplierCategoryRelation) Delete(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelation) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelation) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}
