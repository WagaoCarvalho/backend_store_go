package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleRepo_GetByID(t *testing.T) {
	t.Run("successfully get sale by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)

		clientID := int64(100)
		userID := int64(200)
		saleDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		totalAmount := 150.50
		createdAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{saleID}).
			Return(&mockDb.MockRowWithIDArgs{
				Values: []any{
					saleID,            // id - int64
					&clientID,         // client_id - *int64
					&userID,           // user_id - *int64
					saleDate,          // sale_date - time.Time
					totalAmount,       // total_amount - float64
					float64(10.00),    // total_discount - float64
					"credit",          // payment_type - string
					"completed",       // status - string
					"Test sale notes", // notes - string
					1,                 // version - int
					createdAt,         // created_at - time.Time
					updatedAt,         // updated_at - time.Time
				},
			})

		result, err := repo.GetByID(ctx, saleID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, saleID, result.ID)
		assert.Equal(t, &clientID, result.ClientID)
		assert.Equal(t, &userID, result.UserID)
		assert.Equal(t, saleDate, result.SaleDate)
		assert.Equal(t, totalAmount, result.TotalAmount)
		assert.Equal(t, float64(10.00), result.TotalDiscount)
		assert.Equal(t, "credit", result.PaymentType)
		assert.Equal(t, "completed", result.Status)
		assert.Equal(t, "Test sale notes", result.Notes)
		assert.Equal(t, 1, result.Version)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when sale does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(999)

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{saleID}).
			Return(&mockDb.MockRowWithIDArgs{
				Err: pgx.ErrNoRows,
			})

		result, err := repo.GetByID(ctx, saleID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.NotErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), "não encontrado")

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)
		dbError := errors.New("connection failed")

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{saleID}).
			Return(&mockDb.MockRowWithIDArgs{
				Err: dbError,
			})

		result, err := repo.GetByID(ctx, saleID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())
		assert.Contains(t, err.Error(), "erro ao buscar")

		mockDB.AssertExpectations(t)
	})

	t.Run("handle null client_id and user_id correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(2)

		saleDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)
		totalAmount := 200.00
		createdAt := time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{saleID}).
			Return(&mockDb.MockRowWithIDArgs{
				Values: []any{
					saleID,        // id - int64
					nil,           // client_id - nil
					nil,           // user_id - nil
					saleDate,      // sale_date - time.Time
					totalAmount,   // total_amount - float64
					float64(0.00), // total_discount - float64
					"cash",        // payment_type - string
					"pending",     // status - string
					"",            // notes - string vazio
					1,             // version - int
					createdAt,     // created_at - time.Time
					updatedAt,     // updated_at - time.Time
				},
			})

		result, err := repo.GetByID(ctx, saleID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, saleID, result.ID)
		assert.Nil(t, result.ClientID)
		assert.Nil(t, result.UserID)
		assert.Equal(t, "", result.Notes)
		assert.Equal(t, "pending", result.Status)

		mockDB.AssertExpectations(t)
	})
}

func TestSaleRepo_GetByClientID(t *testing.T) {
	t.Run("return empty list when no sales found for client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(999)
		limit := 10
		offset := 0
		orderBy := "sale_date"
		orderDir := "ASC"

		mockRows := new(mockDb.MockRows)
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return(nil)

		// Execution
		sales, err := repo.GetByClientID(ctx, clientID, limit, offset, orderBy, orderDir)

		// Assertion
		assert.NoError(t, err)
		assert.Empty(t, sales)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

}

func TestSaleRepo_GetByUserID(t *testing.T) {
	t.Run("delegate to listByField with user_id field", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(200)

		// Apenas verifica que a delegação acontece
		mockRows := new(mockDb.MockRows)
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return(nil)

		sales, err := repo.GetByUserID(ctx, userID, 10, 0, "sale_date", "DESC")

		assert.NoError(t, err)
		assert.Empty(t, sales)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSaleRepo_GetByStatus(t *testing.T) {
	t.Run("delegate to listByField with status field", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		status := "completed"

		// Apenas verifica que a delegação acontece
		mockRows := new(mockDb.MockRows)
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return(nil)

		sales, err := repo.GetByStatus(ctx, status, 10, 0, "sale_date", "DESC")

		assert.NoError(t, err)
		assert.Empty(t, sales)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSaleRepo_GetByDateRange(t *testing.T) {
	t.Run("successfully get sales by date range", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
		limit := 10
		offset := 0

		mockRows := new(mockDb.MockRows)

		// CORREÇÃO: Passar os args como slice de interface
		mockDB.On("Query", ctx, mock.Anything, []interface{}{start, end, limit, offset}).
			Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockRows.On("Close").Return(nil)

		sales, err := repo.GetByDateRange(ctx, start, end, limit, offset, "sale_date", "DESC")

		assert.NoError(t, err)
		assert.Len(t, sales, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		dbError := errors.New("connection failed")
		mockDB.On("Query", ctx, mock.Anything, []interface{}{start, end, 10, 0}).
			Return(nil, dbError)

		sales, err := repo.GetByDateRange(ctx, start, end, 10, 0, "sale_date", "DESC")

		assert.Error(t, err)
		assert.Nil(t, sales)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
