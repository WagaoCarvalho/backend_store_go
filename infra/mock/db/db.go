package mock

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	call := m.Called(ctx, query, args)

	result := call.Get(0)
	switch v := result.(type) {
	case *MockRow:
		return v
	case MockRow:
		return &v
	case *MockRowWithID:
		return v
	case MockRowWithID:
		return &v
	default:
		panic("unexpected type returned from mock")
	}
}

func (m *MockDatabase) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	call := m.Called(ctx, query, args)
	return call.Get(0).(pgconn.CommandTag), call.Error(1)
}

func (m *MockDatabase) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	call := m.Called(ctx, query, args)
	return call.Get(0).(pgx.Rows), call.Error(1)
}
