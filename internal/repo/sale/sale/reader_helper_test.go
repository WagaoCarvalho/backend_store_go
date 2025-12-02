package repo

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Função auxiliar para criar uma sale de teste
func createTestSale() *models.Sale {
	return &models.Sale{
		ID:                1,
		ClientID:          utils.Int64Ptr(100),
		UserID:            utils.Int64Ptr(200),
		SaleDate:          *utils.TimePtrFromString("2024-01-15"),
		TotalAmount:       *utils.Float64Ptr(150.50),
		TotalSaleDiscount: 10.00,
		PaymentType:       "credit",
		Status:            "completed",
		Notes:             "Test sale",
		Version:           1,
		CreatedAt:         *utils.TimePtrFromString("2024-01-15T10:00:00Z"),
		UpdatedAt:         *utils.TimePtrFromString("2024-01-15T10:00:00Z"),
	}
}

func TestSaleRepo_ListByField(t *testing.T) {

	t.Run("successfully list sales by field", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		expectedSales := []*models.Sale{createTestSale()}
		mockRows := new(mockDb.MockRows)

		// Configurar expectativas
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				// Simular o scan preenchendo os valores com os tipos CORRETOS
				if ptr, ok := args.Get(0).(*int64); ok {
					*ptr = expectedSales[0].ID
				}
				if ptr, ok := args.Get(1).(*int64); ok {
					if expectedSales[0].ClientID != nil {
						*ptr = *expectedSales[0].ClientID
					}
				}
				if ptr, ok := args.Get(2).(*int64); ok {
					if expectedSales[0].UserID != nil {
						*ptr = *expectedSales[0].UserID
					}
				}
				if ptr, ok := args.Get(3).(*time.Time); ok { // SaleDate é time.Time
					*ptr = expectedSales[0].SaleDate
				}
				if ptr, ok := args.Get(4).(*float64); ok { // TotalAmount é float64
					*ptr = expectedSales[0].TotalAmount
				}
				if ptr, ok := args.Get(5).(*float64); ok { // TotalDiscount é float64 (não ponteiro)
					*ptr = expectedSales[0].TotalSaleDiscount
				}
				if ptr, ok := args.Get(6).(*string); ok { // PaymentType é string
					*ptr = string(expectedSales[0].PaymentType)
				}
				if ptr, ok := args.Get(7).(*string); ok { // Status é string
					*ptr = string(expectedSales[0].Status)
				}
				if ptr, ok := args.Get(8).(*string); ok { // Notes é *string

					*ptr = expectedSales[0].Notes

				}
				if ptr, ok := args.Get(9).(*int); ok { // Version é int64
					*ptr = expectedSales[0].Version
				}
				if ptr, ok := args.Get(10).(*time.Time); ok { // CreatedAt é time.Time
					*ptr = expectedSales[0].CreatedAt
				}
				if ptr, ok := args.Get(11).(*time.Time); ok { // UpdatedAt é time.Time
					*ptr = expectedSales[0].UpdatedAt
				}
			}).Return(nil)
		mockRows.On("Close").Return(nil)

		// Execution
		sales, err := repo.listByField(ctx, "client_id", 100, 10, 0, "created_at", "DESC")

		// Assertion
		assert.NoError(t, err)
		assert.Len(t, sales, 1)
		assert.Equal(t, expectedSales[0].ID, sales[0].ID)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
	t.Run("return empty list when no sales found", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return(nil)

		// Execution
		sales, err := repo.listByField(ctx, "user_id", 999, 10, 0, "sale_date", "ASC")

		// Assertion
		assert.NoError(t, err)
		assert.Empty(t, sales)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return error when database query fails", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		dbError := errors.New("connection failed")
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(nil, dbError)

		// Execution
		sales, err := repo.listByField(ctx, "status", "pending", 10, 0, "total_amount", "DESC")

		// Assertion
		assert.Error(t, err)
		assert.Nil(t, sales)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when scan fails", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		scanError := errors.New("scan failed")

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(scanError)
		mockRows.On("Close").Return(nil)

		// Execution
		sales, err := repo.listByField(ctx, "payment_type", "credit", 5, 0, "id", "ASC")

		// Assertion
		assert.Error(t, err)
		assert.Nil(t, sales)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, scanError.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSanitizeOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid sale_date",
			input:    "sale_date",
			expected: "sale_date",
		},
		{
			name:     "valid total_amount",
			input:    "total_amount",
			expected: "total_amount",
		},
		{
			name:     "empty string returns default",
			input:    "",
			expected: "sale_date",
		},
		{
			name:     "invalid value returns default",
			input:    "invalid_field",
			expected: "sale_date",
		},
		{
			name:     "case sensitive check - uppercase",
			input:    "SALE_DATE",
			expected: "sale_date",
		},
		{
			name:     "case sensitive check - mixed case",
			input:    "Total_Amount",
			expected: "sale_date",
		},
		{
			name:     "sql injection attempt returns default",
			input:    "'; DROP TABLE sales; --",
			expected: "sale_date",
		},
		{
			name:     "whitespace returns default",
			input:    "  sale_date  ",
			expected: "sale_date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeOrderBy(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeOrderDir(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase desc returns DESC",
			input:    "desc",
			expected: "DESC",
		},
		{
			name:     "uppercase DESC returns DESC",
			input:    "DESC",
			expected: "DESC",
		},
		{
			name:     "mixed case DeSc returns DESC",
			input:    "DeSc",
			expected: "DESC",
		},
		{
			name:     "empty string returns ASC",
			input:    "",
			expected: "ASC",
		},
		{
			name:     "invalid value returns ASC",
			input:    "invalid",
			expected: "ASC",
		},
		{
			name:     "asc returns ASC",
			input:    "asc",
			expected: "ASC",
		},
		{
			name:     "ASC returns ASC",
			input:    "ASC",
			expected: "ASC",
		},
		{
			name:     "whitespace with desc returns DESC",
			input:    "  desc  ",
			expected: "DESC",
		},
		{
			name:     "whitespace with DESC returns DESC",
			input:    "  DESC  ",
			expected: "DESC",
		},
		{
			name:     "tab with desc returns DESC",
			input:    "\tdesc\t",
			expected: "DESC",
		},
		{
			name:     "newline with desc returns DESC",
			input:    "\ndesc\n",
			expected: "DESC",
		},
		{
			name:     "sql injection attempt returns ASC",
			input:    "'; DROP TABLE sales; --",
			expected: "ASC",
		},
		{
			name:     "numbers returns ASC",
			input:    "123",
			expected: "ASC",
		},
		{
			name:     "only whitespace returns ASC",
			input:    "   ",
			expected: "ASC",
		},
		{
			name:     "whitespace with asc returns ASC",
			input:    "  asc  ",
			expected: "ASC",
		},
		{
			name:     "multiple spaces returns ASC",
			input:    "     ",
			expected: "ASC",
		},
		{
			name:     "mixed whitespace with desc returns DESC",
			input:    " \t\n desc \t\n ",
			expected: "DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeOrderDir(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Testes de integração para verificar o comportamento em conjunto
func TestSanitizeFunctionsIntegration(t *testing.T) {
	t.Run("valid order by and direction combination", func(t *testing.T) {
		orderBy := sanitizeOrderBy("total_amount")
		orderDir := sanitizeOrderDir("desc")

		assert.Equal(t, "total_amount", orderBy)
		assert.Equal(t, "DESC", orderDir)
	})

	t.Run("invalid order by with valid direction", func(t *testing.T) {
		orderBy := sanitizeOrderBy("invalid_field")
		orderDir := sanitizeOrderDir("DESC")

		assert.Equal(t, "sale_date", orderBy)
		assert.Equal(t, "DESC", orderDir)
	})

	t.Run("valid order by with invalid direction", func(t *testing.T) {
		orderBy := sanitizeOrderBy("sale_date")
		orderDir := sanitizeOrderDir("invalid")

		assert.Equal(t, "sale_date", orderBy)
		assert.Equal(t, "ASC", orderDir)
	})

	t.Run("both invalid returns defaults", func(t *testing.T) {
		orderBy := sanitizeOrderBy("invalid")
		orderDir := sanitizeOrderDir("invalid")

		assert.Equal(t, "sale_date", orderBy)
		assert.Equal(t, "ASC", orderDir)
	})
}

// Testes de performance/boundary
func TestSanitizeOrderDir_Performance(t *testing.T) {
	t.Run("very long string", func(t *testing.T) {
		longString := strings.Repeat("a", 10000)
		result := sanitizeOrderDir(longString)
		assert.Equal(t, "ASC", result)
	})

	t.Run("special characters", func(t *testing.T) {
		specialChars := "!@#$%^&*()"
		result := sanitizeOrderDir(specialChars)
		assert.Equal(t, "ASC", result)
	})
}
