package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQueryBuilder para testes
type MockQueryBuilder struct {
	mock.Mock
}

func (m *MockQueryBuilder) AddCondition(field string, value any) {
	m.Called(field, value)
}

func (m *MockQueryBuilder) AddILIKECondition(field string, value string) {
	m.Called(field, value)
}

// SimpleQueryBuilder - Implementação simples para testes de integração
type SimpleQueryBuilder struct {
	conditions []string
	args       []any
	ilikeCount int
}

func (s *SimpleQueryBuilder) AddCondition(field string, value any) {
	s.conditions = append(s.conditions, field)
	s.args = append(s.args, value)
}

func (s *SimpleQueryBuilder) AddILIKECondition(field string, value string) {
	s.conditions = append(s.conditions, field+" ILIKE")
	s.args = append(s.args, value)
	s.ilikeCount++
}

func TestTextFilter(t *testing.T) {
	t.Run("should apply ILIKE condition when value is not empty", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		filter := TextFilter{
			Field: "product_name",
			Value: "Notebook",
		}

		mockQB.On("AddILIKECondition", "product_name", "Notebook").Once()

		filter.Apply(mockQB)

		mockQB.AssertExpectations(t)
	})

	t.Run("should not apply condition when value is empty", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		filter := TextFilter{
			Field: "product_name",
			Value: "",
		}

		filter.Apply(mockQB)

		mockQB.AssertNotCalled(t, "AddILIKECondition")
		mockQB.AssertNotCalled(t, "AddCondition")
	})

	t.Run("should work with SimpleQueryBuilder implementation", func(t *testing.T) {
		simpleQB := &SimpleQueryBuilder{}
		filter := TextFilter{
			Field: "manufacturer",
			Value: "Dell",
		}

		filter.Apply(simpleQB)

		assert.Len(t, simpleQB.conditions, 1)
		assert.Contains(t, simpleQB.conditions[0], "manufacturer")
		assert.Len(t, simpleQB.args, 1)
		assert.Equal(t, "Dell", simpleQB.args[0])
		assert.Equal(t, 1, simpleQB.ilikeCount)
	})
}

func TestEqualFilter(t *testing.T) {
	t.Run("should apply equal condition for string type", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		value := "active"
		filter := EqualFilter[string]{
			Field: "status",
			Value: &value,
		}

		mockQB.On("AddCondition", "status =", "active").Once()

		filter.Apply(mockQB)

		mockQB.AssertExpectations(t)
	})

	t.Run("should apply equal condition for int type", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		value := int64(123)
		filter := EqualFilter[int64]{
			Field: "supplier_id",
			Value: &value,
		}

		mockQB.On("AddCondition", "supplier_id =", int64(123)).Once()

		filter.Apply(mockQB)

		mockQB.AssertExpectations(t)
	})

	t.Run("should not apply condition when value is nil", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		filter := EqualFilter[string]{
			Field: "status",
			Value: nil,
		}

		filter.Apply(mockQB)

		mockQB.AssertNotCalled(t, "AddCondition")
		mockQB.AssertNotCalled(t, "AddILIKECondition")
	})

	t.Run("should work with SimpleQueryBuilder for int type", func(t *testing.T) {
		simpleQB := &SimpleQueryBuilder{}
		value := 42
		filter := EqualFilter[int]{
			Field: "category_id",
			Value: &value,
		}

		filter.Apply(simpleQB)

		assert.Len(t, simpleQB.conditions, 1)
		assert.Equal(t, "category_id =", simpleQB.conditions[0])
		assert.Len(t, simpleQB.args, 1)
		assert.Equal(t, 42, simpleQB.args[0])
	})
}

func TestRangeFilter(t *testing.T) {
	t.Run("should apply both min and max conditions", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		minVal := 10.0
		maxVal := 100.0
		filter := RangeFilter[float64]{
			FieldMin: "price",
			FieldMax: "price",
			Min:      &minVal,
			Max:      &maxVal,
		}

		mockQB.On("AddCondition", "price >=", 10.0).Once()
		mockQB.On("AddCondition", "price <=", 100.0).Once()

		filter.Apply(mockQB)

		mockQB.AssertExpectations(t)
	})

	t.Run("should apply only min condition when max is nil", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		minVal := 5
		filter := RangeFilter[int]{
			FieldMin: "quantity",
			FieldMax: "quantity",
			Min:      &minVal,
			Max:      nil,
		}

		mockQB.On("AddCondition", "quantity >=", 5).Once()

		filter.Apply(mockQB)

		mockQB.AssertExpectations(t)
		mockQB.AssertNotCalled(t, "AddCondition", "quantity <=", mock.Anything)
	})

	t.Run("should not apply any condition when both min and max are nil", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)
		filter := RangeFilter[float64]{
			FieldMin: "price",
			FieldMax: "price",
			Min:      nil,
			Max:      nil,
		}

		filter.Apply(mockQB)

		mockQB.AssertNotCalled(t, "AddCondition")
		mockQB.AssertNotCalled(t, "AddILIKECondition")
	})

	t.Run("should work with SimpleQueryBuilder for range filters", func(t *testing.T) {
		simpleQB := &SimpleQueryBuilder{}
		minVal := 1.0
		maxVal := 50.0
		filter := RangeFilter[float64]{
			FieldMin: "weight",
			FieldMax: "weight",
			Min:      &minVal,
			Max:      &maxVal,
		}

		filter.Apply(simpleQB)

		assert.Len(t, simpleQB.conditions, 2)
		assert.Contains(t, simpleQB.conditions[0], "weight >=")
		assert.Contains(t, simpleQB.conditions[1], "weight <=")
		assert.Len(t, simpleQB.args, 2)
		assert.Equal(t, 1.0, simpleQB.args[0])
		assert.Equal(t, 50.0, simpleQB.args[1])
	})
}

func TestFilterComposition(t *testing.T) {
	t.Run("should apply multiple filters in sequence", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)

		filters := []FilterCondition{
			TextFilter{Field: "name", Value: "test"},
			EqualFilter[string]{Field: "status", Value: stringPtr("active")},
			RangeFilter[float64]{
				FieldMin: "price",
				FieldMax: "price",
				Min:      float64Ptr(10.0),
				Max:      float64Ptr(100.0),
			},
		}

		mockQB.On("AddILIKECondition", "name", "test").Once()
		mockQB.On("AddCondition", "status =", "active").Once()
		mockQB.On("AddCondition", "price >=", 10.0).Once()
		mockQB.On("AddCondition", "price <=", 100.0).Once()

		for _, filter := range filters {
			filter.Apply(mockQB)
		}

		mockQB.AssertExpectations(t)
	})

	t.Run("should handle mixed nil and non-nil values correctly", func(t *testing.T) {
		mockQB := new(MockQueryBuilder)

		filters := []FilterCondition{
			TextFilter{Field: "name", Value: ""},             // Não deve aplicar
			EqualFilter[string]{Field: "status", Value: nil}, // Não deve aplicar
			RangeFilter[int]{ // Deve aplicar apenas o min
				FieldMin: "quantity",
				FieldMax: "quantity",
				Min:      intPtr(5),
				Max:      nil,
			},
		}

		mockQB.On("AddCondition", "quantity >=", 5).Once()

		for _, filter := range filters {
			filter.Apply(mockQB)
		}

		mockQB.AssertExpectations(t)
		mockQB.AssertNumberOfCalls(t, "AddCondition", 1)
		mockQB.AssertNumberOfCalls(t, "AddILIKECondition", 0)
	})

	t.Run("should work with SimpleQueryBuilder for composition", func(t *testing.T) {
		simpleQB := &SimpleQueryBuilder{}

		filters := []FilterCondition{
			TextFilter{Field: "product_name", Value: "laptop"},
			EqualFilter[int]{Field: "category_id", Value: intPtr(1)},
			RangeFilter[float64]{
				FieldMin: "min_price",
				FieldMax: "max_price",
				Min:      float64Ptr(100.0),
				Max:      float64Ptr(500.0),
			},
		}

		for _, filter := range filters {
			filter.Apply(simpleQB)
		}

		assert.Len(t, simpleQB.conditions, 4)
		assert.Len(t, simpleQB.args, 4)
		assert.Equal(t, 1, simpleQB.ilikeCount)
	})
}

func TestInterfaceCompliance(t *testing.T) {
	t.Run("should verify QueryBuilder interface implementation", func(t *testing.T) {
		// Teste para garantir que o mock satisfaz a interface
		var _ QueryBuilder = &MockQueryBuilder{}
		var _ QueryBuilder = &SimpleQueryBuilder{}
	})

	t.Run("should verify FilterCondition interface implementations", func(t *testing.T) {
		// Teste para garantir que todos os filtros implementam a interface
		var _ FilterCondition = TextFilter{}
		var _ FilterCondition = EqualFilter[string]{}
		var _ FilterCondition = EqualFilter[int]{}
		var _ FilterCondition = EqualFilter[bool]{}
		var _ FilterCondition = RangeFilter[float64]{}
		var _ FilterCondition = RangeFilter[int]{}
		var _ FilterCondition = RangeFilter[string]{}
	})
}

// Helper functions para criar pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}

func int64Ptr(i int64) *int64 {
	return &i
}
