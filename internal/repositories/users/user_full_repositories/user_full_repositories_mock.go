package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserFullRepository struct {
	mock.Mock
}

func (m *MockUserFullRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error) {
	args := m.Called(ctx, tx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserFullRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}
