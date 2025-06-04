package repository

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	"github.com/stretchr/testify/mock"
)

// Mock do repositório
type MockSupplierCategoryRelationRepo struct {
	mock.Mock
}

func (m *MockSupplierCategoryRelationRepo) Create(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, relation)
	if rel, ok := args.Get(0).(*models.SupplierCategoryRelations); ok {
		return rel, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSupplierCategoryRelationRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*models.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationRepo) Update(ctx context.Context, rel *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, rel)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationRepo) Delete(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationRepo) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationRepo) HasSupplierCategoryRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}
