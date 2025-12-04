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

func TestSanitizeField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid client_id", "client_id", "client_id"},
		{"valid user_id", "user_id", "user_id"},
		{"valid status", "status", "status"},
		{"valid payment_type", "payment_type", "payment_type"},
		{"invalid field", "invalid_field", "client_id"},
		{"empty field", "", "client_id"},
		{"case insensitive", "CLIENT_ID", "client_id"}, // Note: como está implementado, isso retornaria "CLIENT_ID"
		{"sql injection attempt", "1; DROP TABLE sales; --", "client_id"},
		{"field with spaces", "client id", "client_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeField(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Função auxiliar para criar uma sale de teste
// Função auxiliar para criar uma sale de teste atualizada
func createTestSale() *models.Sale {
	return &models.Sale{
		ID:                 1,
		ClientID:           utils.Int64Ptr(100),
		UserID:             utils.Int64Ptr(200),
		SaleDate:           *utils.TimePtrFromString("2024-01-15"),
		TotalItemsAmount:   200.00, // NOVO
		TotalItemsDiscount: 20.00,  // NOVO
		TotalSaleDiscount:  10.00,
		TotalAmount:        170.00, // 200 - 20 - 10
		PaymentType:        "credit",
		Status:             "completed",
		Notes:              "Test sale",
		Version:            1,
		CreatedAt:          *utils.TimePtrFromString("2024-01-15T10:00:00Z"),
		UpdatedAt:          *utils.TimePtrFromString("2024-01-15T10:00:00Z"),
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
		// AGORA SÃO 14 CAMPOS NO SCAN (vs 12 antes)
		mockRows.On("Scan",
			mock.Anything, // id (1)
			mock.Anything, // client_id (2)
			mock.Anything, // user_id (3)
			mock.Anything, // sale_date (4)
			mock.Anything, // total_items_amount (5) ← NOVO
			mock.Anything, // total_items_discount (6) ← NOVO
			mock.Anything, // total_sale_discount (7)
			mock.Anything, // total_amount (8)
			mock.Anything, // payment_type (9)
			mock.Anything, // status (10)
			mock.Anything, // notes (11)
			mock.Anything, // version (12)
			mock.Anything, // created_at (13)
			mock.Anything, // updated_at (14)
		).Run(func(args mock.Arguments) {
			// Simular o scan preenchendo os valores com os tipos CORRETOS
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = expectedSales[0].ID
			}
			if ptr, ok := args.Get(1).(**int64); ok {
				// client_id é **int64 (ponteiro para ponteiro)
				if expectedSales[0].ClientID != nil {
					*ptr = expectedSales[0].ClientID
				} else {
					*ptr = nil
				}
			}
			if ptr, ok := args.Get(2).(**int64); ok {
				// user_id é **int64 (ponteiro para ponteiro)
				if expectedSales[0].UserID != nil {
					*ptr = expectedSales[0].UserID
				} else {
					*ptr = nil
				}
			}
			if ptr, ok := args.Get(3).(*time.Time); ok {
				*ptr = expectedSales[0].SaleDate
			}
			if ptr, ok := args.Get(4).(*float64); ok { // total_items_amount (NOVO)
				*ptr = expectedSales[0].TotalItemsAmount
			}
			if ptr, ok := args.Get(5).(*float64); ok { // total_items_discount (NOVO)
				*ptr = expectedSales[0].TotalItemsDiscount
			}
			if ptr, ok := args.Get(6).(*float64); ok { // total_sale_discount
				*ptr = expectedSales[0].TotalSaleDiscount
			}
			if ptr, ok := args.Get(7).(*float64); ok { // total_amount
				*ptr = expectedSales[0].TotalAmount
			}
			if ptr, ok := args.Get(8).(*string); ok { // payment_type
				*ptr = expectedSales[0].PaymentType
			}
			if ptr, ok := args.Get(9).(*string); ok { // status
				*ptr = expectedSales[0].Status
			}
			if ptr, ok := args.Get(10).(*string); ok { // notes
				*ptr = expectedSales[0].Notes
			}
			if ptr, ok := args.Get(11).(*int); ok { // version
				*ptr = expectedSales[0].Version
			}
			if ptr, ok := args.Get(12).(*time.Time); ok { // created_at
				*ptr = expectedSales[0].CreatedAt
			}
			if ptr, ok := args.Get(13).(*time.Time); ok { // updated_at
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
		mockRows.On("Scan",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything,
		).Return(scanError)
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
			name:     "valid created_at",
			input:    "created_at",
			expected: "created_at",
		},
		{
			name:     "valid total_items_amount", // NOVO
			input:    "total_items_amount",
			expected: "total_items_amount",
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
		// ... resto dos testes
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
