package mock

import (
	"github.com/stretchr/testify/mock"
)

// MockQueryBuilderSql para testes
type MockQueryBuilderSql struct {
	mock.Mock
}

func (m *MockQueryBuilderSql) AddCondition(condition string, value any) {
	m.Called(condition, value)
}

func (m *MockQueryBuilderSql) AddILIKECondition(field string, value string) {
	m.Called(field, value)
}

func (m *MockQueryBuilderSql) AddInCondition(field string, values []any) {
	m.Called(field, values)
}

func (m *MockQueryBuilderSql) AddBetweenCondition(field string, value1, value2 any) {
	m.Called(field, value1, value2)
}

func (m *MockQueryBuilderSql) AddIsNullCondition(field string, isNull bool) {
	m.Called(field, isNull)
}

func (m *MockQueryBuilderSql) AddNotEqualCondition(field string, value any) {
	m.Called(field, value)
}

func (m *MockQueryBuilderSql) AddGreaterThanCondition(field string, value any) {
	m.Called(field, value)
}

func (m *MockQueryBuilderSql) AddLessThanCondition(field string, value any) {
	m.Called(field, value)
}

func (m *MockQueryBuilderSql) AddORCondition(condition string, value any) {
	m.Called(condition, value)
}

func (m *MockQueryBuilderSql) AddPagination(limit, offset int) {
	m.Called(limit, offset)
}

func (m *MockQueryBuilderSql) AddOrderBy(field, direction string) {
	m.Called(field, direction)
}

func (m *MockQueryBuilderSql) Build() (string, []any) {
	args := m.Called()
	return args.String(0), args.Get(1).([]any)
}

func (m *MockQueryBuilderSql) GetArgs() []any {
	args := m.Called()
	return args.Get(0).([]any)
}

func (m *MockQueryBuilderSql) GetQuery() string {
	args := m.Called()
	return args.String(0)
}

// Adicione este método ao seu MockQueryBuilderSql no arquivo mock/filter.go
func (m *MockQueryBuilderSql) GetBaseQuery() string {
	args := m.Called()
	return args.String(0)
}
