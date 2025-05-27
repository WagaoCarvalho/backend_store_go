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

func (m *MockSupplierCategoryRelationService) GetBySupplierId(ctx context.Context, supplierID int64) ([]*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) GetByCategoryId(ctx context.Context, categoryID int64) ([]*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) DeleteById(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) Update(ctx context.Context, rel *supplier_category_relations.SupplierCategoryRelations) (*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, rel)
	result := args.Get(0)

	if result == nil {
		return nil, args.Error(1)
	}

	return result.(*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}
