package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Mock para simular a conexão com o banco de dados
type MockPgxPool struct{}

func (m *MockPgxPool) ParseConfig(connString string) (*pgxpool.Config, error) {
	return &pgxpool.Config{}, nil
}

func (m *MockPgxPool) NewWithConfig(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
	return &pgxpool.Pool{}, nil
}
