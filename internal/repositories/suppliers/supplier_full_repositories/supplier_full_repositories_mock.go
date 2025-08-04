package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockSupplierFullRepository struct {
	mock.Mock
}

func (m *MockSupplierFullRepository) CreateTx(ctx context.Context, tx pgx.Tx, supplier *models.Supplier) (*models.Supplier, error) {
	args := m.Called(ctx, tx, supplier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Supplier), args.Error(1)
}

func (m *MockSupplierFullRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}
