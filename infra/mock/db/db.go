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
	case *MockRowWithIDArgs:
		return row
	case MockRowWithIDArgs:
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

func (m *MockDatabase) QueryNoArgs(ctx context.Context, query string) (pgx.Rows, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(pgx.Rows), args.Error(1)
}

type Scanner interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
	Close()
}

func (m *MockRows) QueryWithScanner(ctx context.Context, query string, args ...any) (Scanner, error) {
	argsCalled := m.Called(ctx, query, args)
	if argsCalled.Get(0) != nil {
		return argsCalled.Get(0).(Scanner), argsCalled.Error(1)
	}
	return nil, argsCalled.Error(1)
}
