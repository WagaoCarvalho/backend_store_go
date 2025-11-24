package mock

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockRows struct {
	mock.Mock

	Rows    []*MockRow
	cursor  int
	RowsErr error // usado no modo predefinido
}

func (m *MockRows) Next() bool {
	if len(m.Rows) > 0 {
		if m.cursor < len(m.Rows) {
			m.cursor++
			return true
		}
		return false
	}

	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	if len(m.Rows) > 0 && m.cursor > 0 {
		row := m.Rows[m.cursor-1]
		return row.Scan(dest...)
	}

	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRows) Close() {
	if len(m.Rows) > 0 {
		m.cursor = 0
		return
	}
	m.Called()
}

func (m *MockRows) Conn() *pgx.Conn {
	args := m.Called()
	val := args.Get(0)
	conn, ok := val.(*pgx.Conn)
	if ok {
		return conn
	}
	return nil
}

func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	val := args.Get(0)
	values, ok := val.([]interface{})
	if !ok {
		values = nil
	}
	return values, args.Error(1)
}

func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	val := args.Get(0)
	bytes, ok := val.([][]byte)
	if ok {
		return bytes
	}
	return [][]byte{}
}

func (m *MockRows) Err() error {
	if len(m.Rows) > 0 {
		return m.RowsErr
	}

	args := m.Called()
	if len(args) == 0 {
		return nil
	}
	return args.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	val := args.Get(0)
	tag, ok := val.(pgconn.CommandTag)
	if ok {
		return tag
	}
	return pgconn.CommandTag{}
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	val := args.Get(0)
	fields, ok := val.([]pgconn.FieldDescription)
	if ok {
		return fields
	}
	return []pgconn.FieldDescription{}
}
