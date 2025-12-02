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
			ClientID:           utils.Int64Ptr(100),
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   120.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        115.00,
			PaymentType:        "credit",
			Status:             "pending",
			Notes:              "Test sale notes",
		}

		// Note: A query retorna apenas id, created_at, updated_at
		// Se houver mais campos, ajuste conforme necessário
		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				int64(1),     // id
				expectedTime, // created_at
				expectedTime, // updated_at
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{
				sale.ClientID,
				sale.UserID,
				sale.SaleDate,
				sale.TotalItemsAmount,
				sale.TotalItemsDiscount,
				sale.TotalSaleDiscount,
				sale.TotalAmount,
				sale.PaymentType,
				sale.Status,
				sale.Notes,
			}).Return(mockRow)

		result, err := repo.Create(ctx, sale)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), sale.ID)
		assert.Equal(t, expectedTime, sale.CreatedAt)
		assert.Equal(t, expectedTime, sale.UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(100),
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   120.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        115.00,
			PaymentType:        "credit",
			Status:             "pending",
			Notes:              "Test sale notes",
		}

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

	t.Run("return ErrDuplicate when unique violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(100),
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   120.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        115.00,
			PaymentType:        "credit",
			Status:             "pending",
			Notes:              "Test sale notes",
		}

		uniqueError := &pgconn.PgError{
			Code:           "23505",
			Message:        "violação de unicidade",
			ConstraintName: "sales_unique_constraint",
		}

		mockRow := &mockDb.MockRow{Err: uniqueError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		result, err := repo.Create(ctx, sale)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, errMsg.ErrDuplicate.Error())
		assert.ErrorContains(t, err, "sales_unique_constraint")

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(100),
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   120.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        115.00,
			PaymentType:        "credit",
			Status:             "pending",
			Notes:              "Test sale notes",
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

	t.Run("handle nil pointer fields gracefully", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		expectedTime := time.Now()

		sale := &models.Sale{
			// ClientID e UserID são nil
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   120.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        115.00,
			PaymentType:        "credit",
			Status:             "pending",
			Notes:              "Test sale notes",
		}

		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				int64(1),     // id
				expectedTime, // created_at
				expectedTime, // updated_at
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{
				sale.ClientID, // Será nil
				sale.UserID,   // Será nil
				sale.SaleDate,
				sale.TotalItemsAmount,
				sale.TotalItemsDiscount,
				sale.TotalSaleDiscount,
				sale.TotalAmount,
				sale.PaymentType,
				sale.Status,
				sale.Notes,
			}).Return(mockRow)

		result, err := repo.Create(ctx, sale)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), sale.ID)

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
			ID:                 1,
			ClientID:           utils.Int64Ptr(100),
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   200.00,
			TotalItemsDiscount: 15.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        180.00,
			PaymentType:        "credit",
			Status:             "completed",
			Notes:              "Updated sale notes",
			Version:            1,
		}

		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				expectedTime, // updated_at
				2,            // version incrementado
			},
		}

		// CORREÇÃO: A query tem 12 argumentos ($1 a $12)
		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{
				sale.ClientID,           // $1
				sale.UserID,             // $2
				sale.SaleDate,           // $3
				sale.TotalItemsAmount,   // $4
				sale.TotalItemsDiscount, // $5
				sale.TotalSaleDiscount,  // $6 - ESTE ESTAVA FALTANDO NO TESTE ORIGINAL
				sale.TotalAmount,        // $7 - ESTE ESTAVA FALTANDO NO TESTE ORIGINAL
				sale.PaymentType,        // $8
				sale.Status,             // $9
				sale.Notes,              // $10
				sale.ID,                 // $11
				sale.Version,            // $12
			}).
			Return(mockRow)

		err := repo.Update(ctx, sale)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, sale.UpdatedAt)
		assert.Equal(t, 2, sale.Version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when sale does not exist or version mismatch", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()

		sale := &models.Sale{
			ID:                 999,
			ClientID:           utils.Int64Ptr(100),
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   200.00,
			TotalItemsDiscount: 15.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        180.00,
			PaymentType:        "credit",
			Status:             "completed",
			Notes:              "Updated sale notes",
			Version:            1,
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		// Usando mock.Anything para os argumentos pois só queremos testar o erro
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
			ID:                 1,
			ClientID:           utils.Int64Ptr(999), // FK inválido
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   200.00,
			TotalItemsDiscount: 15.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        180.00,
			PaymentType:        "credit",
			Status:             "completed",
			Notes:              "Updated sale notes",
			Version:            1,
		}

		fkError := &pgconn.PgError{Code: "23503", Message: "violação de chave estrangeira"}

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
			ID:                 1,
			ClientID:           utils.Int64Ptr(100),
			UserID:             utils.Int64Ptr(200),
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   200.00,
			TotalItemsDiscount: 15.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        180.00,
			PaymentType:        "credit",
			Status:             "completed",
			Notes:              "Updated sale notes",
			Version:            1,
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

	t.Run("handle nil pointer fields gracefully", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		expectedTime := time.Now()

		sale := &models.Sale{
			ID: 1,
			// ClientID e UserID são nil
			SaleDate:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalItemsAmount:   200.00,
			TotalItemsDiscount: 15.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        180.00,
			PaymentType:        "credit",
			Status:             "completed",
			Notes:              "Updated sale notes",
			Version:            1,
		}

		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				expectedTime, // updated_at
				2,            // version incrementado
			},
		}

		// Teste com campos nil
		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{
				sale.ClientID,           // $1 - será nil
				sale.UserID,             // $2 - será nil
				sale.SaleDate,           // $3
				sale.TotalItemsAmount,   // $4
				sale.TotalItemsDiscount, // $5
				sale.TotalSaleDiscount,  // $6
				sale.TotalAmount,        // $7
				sale.PaymentType,        // $8
				sale.Status,             // $9
				sale.Notes,              // $10
				sale.ID,                 // $11
				sale.Version,            // $12
			}).
			Return(mockRow)

		err := repo.Update(ctx, sale)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, sale.UpdatedAt)
		assert.Equal(t, 2, sale.Version)
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
