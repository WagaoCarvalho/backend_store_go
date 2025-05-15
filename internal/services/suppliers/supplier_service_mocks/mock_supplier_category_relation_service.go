package supplier_service_mocks

import (
	"context"

	supplier_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	"github.com/stretchr/testify/mock"
)

type MockSupplierCategoryRelationService struct {
	mock.Mock
}

func (m *MockSupplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Get(0).(*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) GetBySupplier(ctx context.Context, supplierID int64) ([]*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) GetByCategory(ctx context.Context, categoryID int64) ([]*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) Delete(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) DeleteAll(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}
