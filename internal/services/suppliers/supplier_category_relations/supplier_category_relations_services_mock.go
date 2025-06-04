package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	"github.com/stretchr/testify/mock"
)

type MockSupplierCategoryRelationService struct {
	mock.Mock
}

func (m *MockSupplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID, categoryID)
	var rel *models.SupplierCategoryRelations
	if tmp := args.Get(0); tmp != nil {
		rel = tmp.(*models.SupplierCategoryRelations)
	}
	return rel, args.Error(1)
}

func (m *MockSupplierCategoryRelationService) GetBySupplierId(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID)
	var rels []*models.SupplierCategoryRelations
	if tmp := args.Get(0); tmp != nil {
		rels = tmp.([]*models.SupplierCategoryRelations)
	}
	return rels, args.Error(1)
}

func (m *MockSupplierCategoryRelationService) GetByCategoryId(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	var rels []*models.SupplierCategoryRelations
	if tmp := args.Get(0); tmp != nil {
		rels = tmp.([]*models.SupplierCategoryRelations)
	}
	return rels, args.Error(1)
}

func (m *MockSupplierCategoryRelationService) Update(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, relation)
	var rel *models.SupplierCategoryRelations
	if tmp := args.Get(0); tmp != nil {
		rel = tmp.(*models.SupplierCategoryRelations)
	}
	return rel, args.Error(1)
}

func (m *MockSupplierCategoryRelationService) DeleteById(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}
