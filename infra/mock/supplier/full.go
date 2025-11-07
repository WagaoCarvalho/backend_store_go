package mock

import (
	"context"

	model_full "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/full"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockSupplierFull struct {
	mock.Mock
}

func (m *MockSupplierFull) CreateFull(ctx context.Context, supplier *model_full.SupplierFull) (*model_full.SupplierFull, error) {
	args := m.Called(ctx, supplier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model_full.SupplierFull), args.Error(1)
}

func (m *MockSupplierFull) CreateTx(ctx context.Context, tx pgx.Tx, supplier *models.Supplier) (*models.Supplier, error) {
	args := m.Called(ctx, tx, supplier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Supplier), args.Error(1)
}

func (m *MockSupplierFull) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}
