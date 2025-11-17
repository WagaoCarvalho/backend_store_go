package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleRepo_Create(t *testing.T) {
	t.Run("successfully create sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		expectedTime := time.Now()

		sale := &models.Sale{
			ClientID:      utils.Int64Ptr(100),
			UserID:        utils.Int64Ptr(200),
			SaleDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalAmount:   150.50,
			TotalDiscount: 10.00,
			PaymentType:   "credit",
			Status:        "pending",
			Notes:         "Test sale notes",
		}

		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				int64(1),     // id
				1,            // version
				expectedTime, // created_at
				expectedTime, // updated_at
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{
				sale.ClientID,
				sale.UserID,
				sale.SaleDate,
				sale.TotalAmount,
				sale.TotalDiscount,
				sale.PaymentType,
				sale.Status,
				sale.Notes,
			}).
			Return(mockRow)

		result, err := repo.Create(ctx, sale)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, int64(100), *result.ClientID)
		assert.Equal(t, int64(200), *result.UserID)
		assert.Equal(t, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), result.SaleDate)
		assert.Equal(t, 150.50, result.TotalAmount)
		assert.Equal(t, 10.00, result.TotalDiscount)
		assert.Equal(t, "credit", result.PaymentType)
		assert.Equal(t, "pending", result.Status)
		assert.Equal(t, "Test sale notes", result.Notes)
		assert.Equal(t, 1, result.Version)
		assert.Equal(t, expectedTime, result.CreatedAt)
		assert.Equal(t, expectedTime, result.UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ClientID:      utils.Int64Ptr(100),
			UserID:        utils.Int64Ptr(200),
			SaleDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalAmount:   150.50,
			TotalDiscount: 10.00,
			PaymentType:   "credit",
			Status:        "pending",
			Notes:         "Test sale notes",
		}

		// Simulando erro de violação de chave estrangeira
		fkError := &pgconn.PgError{
			Code:    "23503",
			Message: "violação de chave estrangeira",
		}

		mockRow := &mockDb.MockRow{Err: fkError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		result, err := repo.Create(ctx, sale)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ClientID:      utils.Int64Ptr(100),
			UserID:        utils.Int64Ptr(200),
			SaleDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalAmount:   150.50,
			TotalDiscount: 10.00,
			PaymentType:   "credit",
			Status:        "pending",
			Notes:         "Test sale notes",
		}

		dbError := errors.New("connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		result, err := repo.Create(ctx, sale)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})
}

func TestSaleRepo_Update(t *testing.T) {
	t.Run("successfully update sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		expectedTime := time.Now()

		sale := &models.Sale{
			ID:            int64(1),
			ClientID:      utils.Int64Ptr(100),
			UserID:        utils.Int64Ptr(200),
			SaleDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalAmount:   150.50,
			TotalDiscount: 10.00,
			PaymentType:   "credit",
			Status:        "completed",
			Notes:         "Updated sale notes",
			Version:       1,
		}

		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				expectedTime, // updated_at
				2,            // version (incrementado)
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{
				sale.ClientID,
				sale.UserID,
				sale.SaleDate,
				sale.TotalAmount,
				sale.TotalDiscount,
				sale.PaymentType,
				sale.Status,
				sale.Notes,
				sale.ID,
				sale.Version,
			}).
			Return(mockRow)

		err := repo.Update(ctx, sale)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, sale.UpdatedAt)
		assert.Equal(t, 2, sale.Version) // Version foi incrementado
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when sale does not exist or version mismatch", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ID:            int64(999),
			ClientID:      utils.Int64Ptr(100),
			UserID:        utils.Int64Ptr(200),
			SaleDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalAmount:   150.50,
			TotalDiscount: 10.00,
			PaymentType:   "credit",
			Status:        "completed",
			Notes:         "Updated sale notes",
			Version:       1,
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		err := repo.Update(ctx, sale)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ID:            int64(1),
			ClientID:      utils.Int64Ptr(999), // ClientID inválido
			UserID:        utils.Int64Ptr(200),
			SaleDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalAmount:   150.50,
			TotalDiscount: 10.00,
			PaymentType:   "credit",
			Status:        "completed",
			Notes:         "Updated sale notes",
			Version:       1,
		}

		fkError := &pgconn.PgError{
			Code:    "23503",
			Message: "violação de chave estrangeira",
		}

		mockRow := &mockDb.MockRow{Err: fkError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		err := repo.Update(ctx, sale)

		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when other database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ID:            int64(1),
			ClientID:      utils.Int64Ptr(100),
			UserID:        utils.Int64Ptr(200),
			SaleDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalAmount:   150.50,
			TotalDiscount: 10.00,
			PaymentType:   "credit",
			Status:        "completed",
			Notes:         "Updated sale notes",
			Version:       1,
		}

		dbError := errors.New("connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		err := repo.Update(ctx, sale)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSaleRepo_Delete(t *testing.T) {
	t.Run("successfully delete sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 1,
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.Delete(ctx, saleID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when sale does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(999)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 0,
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.Delete(ctx, saleID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)
		dbError := errors.New("connection failed")

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(nil, dbError)

		err := repo.Delete(ctx, saleID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
