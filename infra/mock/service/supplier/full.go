package services

import (
	"context"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/full"
	"github.com/stretchr/testify/mock"
)

// Mock do serviço de usuário
type MockSupplierFullService struct {
	mock.Mock
}

func (m *MockSupplierFullService) CreateFull(ctx context.Context, supplier *model.SupplierFull) (*model.SupplierFull, error) {
	args := m.Called(ctx, supplier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SupplierFull), args.Error(1)
}
