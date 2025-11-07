package mock

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Para satisfazer completamente a interface pgx.Tx, adicione os métodos que seu código usar:
func (m *MockTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	callArgs := m.Called(append([]any{ctx, sql}, args...)...)
	cmdTag, _ := callArgs.Get(0).(pgconn.CommandTag)
	return cmdTag, callArgs.Error(1)
}

// Adicione outros métodos se necessário
func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	callArgs := m.Called(append([]any{ctx, sql}, args...)...)
	rows, _ := callArgs.Get(0).(pgx.Rows)
	return rows, callArgs.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	callArgs := m.Called(append([]any{ctx, sql}, args...)...)
	row, _ := callArgs.Get(0).(pgx.Row)
	return row
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	br, _ := args.Get(0).(pgx.BatchResults)
	return br
}

func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	conn, _ := args.Get(0).(*pgx.Conn)
	return conn
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columns []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columns, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	return args.Get(0).(pgx.LargeObjects)
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	stmtDesc, _ := args.Get(0).(*pgconn.StatementDescription)
	return stmtDesc, args.Error(1)
}
