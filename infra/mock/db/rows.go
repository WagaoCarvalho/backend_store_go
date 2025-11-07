package mock

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRows) Conn() *pgx.Conn {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*pgx.Conn)
	}
	return nil
}

func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]interface{}), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([][]byte)
	}
	return [][]byte{}
}

func (m *MockRows) Err() error {
	args := m.Called()
	if len(args) == 0 {
		return nil
	}
	return args.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(pgconn.CommandTag)
	}
	return pgconn.CommandTag{}
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]pgconn.FieldDescription)
	}
	return []pgconn.FieldDescription{}
}

func (m *MockRows) Close() {
	m.Called()
}
