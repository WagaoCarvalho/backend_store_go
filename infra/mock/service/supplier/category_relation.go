package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	"github.com/stretchr/testify/mock"
)

type MockSupplierCategoryRelationService struct {
	mock.Mock
}

func (m *MockSupplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelation, bool, error) {
	args := m.Called(ctx, supplierID, categoryID)

	var rel *models.SupplierCategoryRelation
	if tmp := args.Get(0); tmp != nil {
		rel = tmp.(*models.SupplierCategoryRelation)
	}

	return rel, args.Bool(1), args.Error(2)
}

func (m *MockSupplierCategoryRelationService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error) {
	args := m.Called(ctx, supplierID)

	var rels []*models.SupplierCategoryRelation
	if tmp := args.Get(0); tmp != nil {
		rels = tmp.([]*models.SupplierCategoryRelation)
	}

	return rels, args.Error(1)
}

func (m *MockSupplierCategoryRelationService) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error) {
	args := m.Called(ctx, categoryID)

	var rels []*models.SupplierCategoryRelation
	if tmp := args.Get(0); tmp != nil {
		rels = tmp.([]*models.SupplierCategoryRelation)
	}

	return rels, args.Error(1)
}

func (m *MockSupplierCategoryRelationService) DeleteByID(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}
