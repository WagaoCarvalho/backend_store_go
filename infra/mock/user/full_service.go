package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/full"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserFullService struct {
	mock.Mock
}

func (m *MockUserFullService) CreateFull(ctx context.Context, user *models.UserFull) (*models.UserFull, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserFull), args.Error(1)
}

func (m *MockUserFullService) CreateTx(ctx context.Context, tx pgx.Tx, user *models.UserFull) (*models.UserFull, error) {
	args := m.Called(ctx, tx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserFull), args.Error(1)
}

func (m *MockUserFullService) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}
