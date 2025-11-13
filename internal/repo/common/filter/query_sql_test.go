package repo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSQLQueryBuilder(t *testing.T) {
	t.Run("should create SQLQueryBuilder with correct initial state", func(t *testing.T) {
		tableName := "products"
		columns := []string{"id", "name", "price"}
		orderBy := "created_at DESC"

		qb := NewSQLQueryBuilder(tableName, columns, orderBy)

		assert.Equal(t, tableName, qb.tableName)
		assert.Equal(t, columns, qb.columns)
		assert.Equal(t, orderBy, qb.orderBy)
		assert.Equal(t, 1, qb.pos)
		assert.Empty(t, qb.args)
		assert.Contains(t, qb.where.String(), "WHERE 1=1")
	})
}

func TestSQLQueryBuilder_AddCondition(t *testing.T) {
	t.Run("should add condition with proper parameter numbering", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "id")

		qb.AddCondition("status =", "active")
		qb.AddCondition("price >=", 100.0)

		query, args := qb.Build(10, 0)

		assert.Contains(t, query, "status = $1")
		assert.Contains(t, query, "price >= $2")
		assert.Len(t, args, 2)
		assert.Equal(t, "active", args[0])
		assert.Equal(t, 100.0, args[1])
	})

	t.Run("should handle multiple data types in conditions", func(t *testing.T) {
		qb := NewSQLQueryBuilder("users", []string{"id", "name"}, "name")

		qb.AddCondition("age >=", 18)
		qb.AddCondition("is_active =", true)
		qb.AddCondition("name =", "John")
		qb.AddCondition("salary <=", 5000.50)

		_, args := qb.Build(10, 0)

		assert.Len(t, args, 4)
		assert.Equal(t, 18, args[0])
		assert.Equal(t, true, args[1])
		assert.Equal(t, "John", args[2])
		assert.Equal(t, 5000.50, args[3])
	})

	t.Run("should increment parameter position correctly", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "id")

		// Add multiple conditions
		for i := 1; i <= 5; i++ {
			qb.AddCondition(fmt.Sprintf("field%d =", i), i)
		}

		query, args := qb.Build(10, 0)

		for i := 1; i <= 5; i++ {
			assert.Contains(t, query, fmt.Sprintf("field%d = $%d", i, i))
		}
		assert.Len(t, args, 5)
		assert.Equal(t, 6, qb.pos) // Next position should be 6
	})
}

func TestSQLQueryBuilder_AddILIKECondition(t *testing.T) {
	t.Run("should add ILIKE condition with proper formatting", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "name")

		qb.AddILIKECondition("name", "laptop")
		qb.AddILIKECondition("description", "gaming")

		query, args := qb.Build(10, 0)

		assert.Contains(t, query, "name ILIKE '%' || $1 || '%'")
		assert.Contains(t, query, "description ILIKE '%' || $2 || '%'")
		assert.Len(t, args, 2)
		assert.Equal(t, "laptop", args[0])
		assert.Equal(t, "gaming", args[1])
	})

	t.Run("should mix ILIKE and regular conditions", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "name")

		qb.AddILIKECondition("name", "apple")
		qb.AddCondition("price >=", 1000)
		qb.AddILIKECondition("category", "electronics")

		query, args := qb.Build(10, 0)

		assert.Contains(t, query, "name ILIKE '%' || $1 || '%'")
		assert.Contains(t, query, "price >= $2")
		assert.Contains(t, query, "category ILIKE '%' || $3 || '%'")
		assert.Len(t, args, 3)
		assert.Equal(t, "apple", args[0])
		assert.Equal(t, 1000, args[1])
		assert.Equal(t, "electronics", args[2])
	})
}

func TestSQLQueryBuilder_Build(t *testing.T) {
	t.Run("should build complete query with all components", func(t *testing.T) {
		qb := NewSQLQueryBuilder(
			"products",
			[]string{"id", "name", "price", "category"},
			"created_at DESC, name ASC",
		)

		qb.AddCondition("price >=", 100)
		qb.AddILIKECondition("name", "phone")

		query, args := qb.Build(25, 50)

		// Verify SELECT clause
		assert.Contains(t, query, "SELECT id, name, price, category")

		// Verify FROM clause
		assert.Contains(t, query, "FROM products")

		// Verify WHERE clause
		assert.Contains(t, query, "WHERE 1=1")
		assert.Contains(t, query, "price >= $1")
		assert.Contains(t, query, "name ILIKE '%' || $2 || '%'")

		// Verify ORDER BY
		assert.Contains(t, query, "ORDER BY created_at DESC, name ASC")

		// Verify LIMIT and OFFSET
		assert.Contains(t, query, "LIMIT 25 OFFSET 50")

		// Verify args
		assert.Len(t, args, 2)
		assert.Equal(t, 100, args[0])
		assert.Equal(t, "phone", args[1])
	})

	t.Run("should build query without additional conditions", func(t *testing.T) {
		qb := NewSQLQueryBuilder("users", []string{"*"}, "id")

		query, args := qb.Build(10, 0)

		assert.Contains(t, query, "SELECT *")
		assert.Contains(t, query, "FROM users")
		assert.Contains(t, query, "WHERE 1=1")
		assert.Contains(t, query, "ORDER BY id")
		assert.Contains(t, query, "LIMIT 10 OFFSET 0")
		assert.Empty(t, args)
	})

	t.Run("should handle different limit and offset values", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "id")

		testCases := []struct {
			limit  int
			offset int
		}{
			{10, 0},
			{0, 0},    // zero limit
			{100, 50}, // large values
			{1, 999},  // edge cases
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("limit=%d_offset=%d", tc.limit, tc.offset), func(t *testing.T) {
				query, _ := qb.Build(tc.limit, tc.offset)
				assert.Contains(t, query, fmt.Sprintf("LIMIT %d OFFSET %d", tc.limit, tc.offset))
			})
		}
	})

	t.Run("should trim spaces and format query correctly", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"id", "name"}, "name")

		query, _ := qb.Build(10, 0)

		// Should not have leading/trailing whitespace
		assert.Equal(t, strings.TrimSpace(query), query)

		// Should have proper newlines for readability
		assert.Contains(t, query, "\n")
	})
}

func TestSQLQueryBuilder_IntegrationWithFilters(t *testing.T) {
	t.Run("should work with TextFilter", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "name")
		filter := TextFilter{
			Field: "product_name",
			Value: "laptop",
		}

		filter.Apply(qb)

		query, args := qb.Build(10, 0)

		assert.Contains(t, query, "product_name ILIKE '%' || $1 || '%'")
		assert.Len(t, args, 1)
		assert.Equal(t, "laptop", args[0])
	})

	t.Run("should work with EqualFilter", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "name")
		status := "active"
		filter := EqualFilter[string]{
			Field: "status",
			Value: &status,
		}

		filter.Apply(qb)

		query, args := qb.Build(10, 0)

		assert.Contains(t, query, "status = $1")
		assert.Len(t, args, 1)
		assert.Equal(t, "active", args[0])
	})

	t.Run("should work with RangeFilter", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "price")
		minPrice := 100.0
		maxPrice := 500.0
		filter := RangeFilter[float64]{
			FieldMin: "price",
			FieldMax: "price",
			Min:      &minPrice,
			Max:      &maxPrice,
		}

		filter.Apply(qb)

		query, args := qb.Build(10, 0)

		assert.Contains(t, query, "price >= $1")
		assert.Contains(t, query, "price <= $2")
		assert.Len(t, args, 2)
		assert.Equal(t, 100.0, args[0])
		assert.Equal(t, 500.0, args[1])
	})

	t.Run("should work with multiple filter types", func(t *testing.T) {
		qb := NewSQLQueryBuilder("products", []string{"*"}, "created_at DESC")

		filters := []FilterCondition{
			TextFilter{Field: "name", Value: "gaming"},
			EqualFilter[string]{Field: "category", Value: stringPtr("electronics")},
			RangeFilter[float64]{
				FieldMin: "price",
				FieldMax: "price",
				Min:      float64Ptr(100.0),
				Max:      float64Ptr(1000.0),
			},
		}

		for _, filter := range filters {
			filter.Apply(qb)
		}

		query, args := qb.Build(20, 40)

		assert.Contains(t, query, "name ILIKE '%' || $1 || '%'")
		assert.Contains(t, query, "category = $2")
		assert.Contains(t, query, "price >= $3")
		assert.Contains(t, query, "price <= $4")
		assert.Contains(t, query, "LIMIT 20 OFFSET 40")
		assert.Len(t, args, 4)
	})
}

func TestSQLQueryBuilder_QueryBuilderInterface(t *testing.T) {
	t.Run("should implement QueryBuilder interface", func(t *testing.T) {
		var _ QueryBuilder = &SQLQueryBuilder{}
	})

	t.Run("should be usable polymorphically through QueryBuilder interface", func(t *testing.T) {
		var qb QueryBuilder = NewSQLQueryBuilder("products", []string{"*"}, "id")

		// Should be able to call interface methods
		qb.AddCondition("status =", "active")
		qb.AddILIKECondition("name", "test")

		// Type assertion should work
		sqlQB, ok := qb.(*SQLQueryBuilder)
		assert.True(t, ok)
		assert.Equal(t, "products", sqlQB.tableName)
	})
}
