package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_contact_relations"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockSupplierContactRelationRepo struct {
	mock.Mock
}

func (m *MockSupplierContactRelationRepo) Create(ctx context.Context, relation *models.SupplierContactRelations) (*models.SupplierContactRelations, error) {
	args := m.Called(ctx, relation)
	var result *models.SupplierContactRelations
	if args.Get(0) != nil {
		result = args.Get(0).(*models.SupplierContactRelations)
	}
	return result, args.Error(1)
}

func (m *MockSupplierContactRelationRepo) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierContactRelations) (*models.SupplierContactRelations, error) {
	args := m.Called(ctx, tx, relation)
	var result *models.SupplierContactRelations
	if args.Get(0) != nil {
		result = args.Get(0).(*models.SupplierContactRelations)
	}
	return result, args.Error(1)
}

func (m *MockSupplierContactRelationRepo) HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error) {
	args := m.Called(ctx, supplierID, contactID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSupplierContactRelationRepo) GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelations, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SupplierContactRelations), args.Error(1)
}

func (m *MockSupplierContactRelationRepo) Delete(ctx context.Context, supplierID, contactID int64) error {
	args := m.Called(ctx, supplierID, contactID)
	return args.Error(0)
}

func (m *MockSupplierContactRelationRepo) DeleteAll(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}
