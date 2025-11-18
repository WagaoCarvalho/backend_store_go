package mock

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockDBTransactor struct {
	mock.Mock
}

func (m *MockDBTransactor) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	args := m.Called(ctx, txOptions)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockDBTransactor) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	call := m.Called(ctx, query, args)
	return call.Get(0).(pgx.Row)
}

func (m *MockDBTransactor) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	call := m.Called(ctx, query, args)
	return call.Get(0).(pgconn.CommandTag), call.Error(1)
}

func (m *MockDBTransactor) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	call := m.Called(ctx, query, args)
	return call.Get(0).(pgx.Rows), call.Error(1)
}
