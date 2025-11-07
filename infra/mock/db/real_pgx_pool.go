package mock

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type MockRealPgxPool struct {
	mock.Mock
}

func (m *MockRealPgxPool) ParseConfig(connString string) (*pgxpool.Config, error) {
	args := m.Called(connString)
	return args.Get(0).(*pgxpool.Config), args.Error(1)
}

func (m *MockRealPgxPool) NewWithConfig(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(*pgxpool.Pool), args.Error(1)
}

func (m *MockRealPgxPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
