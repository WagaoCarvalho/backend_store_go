package mock

import (
	"context"
	"fmt"

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

	switch row := result.(type) {
	case *MockRow:
		return row
	case MockRow:
		return &row
	case *MockRowWithInt:
		return row
	case MockRowWithInt:
		return &row
	case *MockRowWithID:
		return row
	case MockRowWithID:
		return &row
	default:
		panic("unexpected type returned from mock")
	}
}

func (m *MockDatabase) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	call := m.Called(ctx, query, args)

	result := call.Get(0)
	if result == nil {
		return pgconn.CommandTag{}, call.Error(1)
	}

	if cmdTag, ok := result.(MockCommandTag); ok {

		return pgconn.NewCommandTag(fmt.Sprintf("UPDATE %d", cmdTag.RowsAffectedCount)), call.Error(1)
	}

	return result.(pgconn.CommandTag), call.Error(1)
}

func (m *MockDatabase) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	call := m.Called(ctx, query, args)
	if v := call.Get(0); v != nil {
		return v.(pgx.Rows), call.Error(1)
	}
	return nil, call.Error(1)
}
