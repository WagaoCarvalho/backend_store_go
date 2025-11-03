package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockSupplierContactRelation struct {
	mock.Mock
}

func (m *MockSupplierContactRelation) Create(ctx context.Context, relation *models.SupplierContactRelation) (*models.SupplierContactRelation, error) {
	args := m.Called(ctx, relation)
	var result *models.SupplierContactRelation
	if args.Get(0) != nil {
		result = args.Get(0).(*models.SupplierContactRelation)
	}
	return result, args.Error(1)
}

func (m *MockSupplierContactRelation) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierContactRelation) (*models.SupplierContactRelation, error) {
	args := m.Called(ctx, tx, relation)
	var result *models.SupplierContactRelation
	if args.Get(0) != nil {
		result = args.Get(0).(*models.SupplierContactRelation)
	}
	return result, args.Error(1)
}

func (m *MockSupplierContactRelation) HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error) {
	args := m.Called(ctx, supplierID, contactID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSupplierContactRelation) GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelation, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SupplierContactRelation), args.Error(1)
}

func (m *MockSupplierContactRelation) Delete(ctx context.Context, supplierID, contactID int64) error {
	args := m.Called(ctx, supplierID, contactID)
	return args.Error(0)
}

func (m *MockSupplierContactRelation) DeleteAll(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}
