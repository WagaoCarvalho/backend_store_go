package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_contact_relations"
	"github.com/stretchr/testify/mock"
)

type MockSupplierContactRelationService struct {
	mock.Mock
}

func (m *MockSupplierContactRelationService) Create(ctx context.Context, supplierID, contactID int64) (*models.SupplierContactRelations, bool, error) {
	args := m.Called(ctx, supplierID, contactID)
	var result *models.SupplierContactRelations
	if args.Get(0) != nil {
		result = args.Get(0).(*models.SupplierContactRelations)
	}
	return result, args.Bool(1), args.Error(2)
}

func (m *MockSupplierContactRelationService) GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelations, error) {
	args := m.Called(ctx, supplierID)
	var result []*models.SupplierContactRelations
	if args.Get(0) != nil {
		result = args.Get(0).([]*models.SupplierContactRelations)
	}
	return result, args.Error(1)
}

func (m *MockSupplierContactRelationService) HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error) {
	args := m.Called(ctx, supplierID, contactID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSupplierContactRelationService) Delete(ctx context.Context, supplierID, contactID int64) error {
	args := m.Called(ctx, supplierID, contactID)
	return args.Error(0)
}

func (m *MockSupplierContactRelationService) DeleteAll(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}
