package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type MockPgxPool struct {
	mock.Mock
}

func (m *MockPgxPool) ParseConfig(connString string) (*pgxpool.Config, error) {
	args := m.Called(connString)
	return args.Get(0).(*pgxpool.Config), args.Error(1)
}

func (m *MockPgxPool) NewWithConfig(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(*pgxpool.Pool), args.Error(1)
}
