package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItemSale_Create(t *testing.T) {
	t.Run("successfully create item sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		expectedTime := time.Now()

		item := &models.SaleItem{
			SaleID:      int64(10),
			ProductID:   int64(20),
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    50.0,
			Description: "desc test",
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
				item.SaleID,
				item.ProductID,
				item.Quantity,
				item.UnitPrice,
				item.Discount,
				item.Tax,
				item.Subtotal,
				item.Description,
			}).
			Return(mockRow)

		result, err := repo.Create(ctx, item)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, int64(10), result.SaleID)
		assert.Equal(t, int64(20), result.ProductID)
		assert.Equal(t, 5, result.Quantity)
		assert.Equal(t, 10.5, result.UnitPrice)
		assert.Equal(t, 1.0, result.Discount)
		assert.Equal(t, 0.5, result.Tax)
		assert.Equal(t, 50.0, result.Subtotal)
		assert.Equal(t, "desc test", result.Description)
		assert.Equal(t, expectedTime, result.CreatedAt)
		assert.Equal(t, expectedTime, result.UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()

		item := &models.SaleItem{
			SaleID:      int64(10),
			ProductID:   int64(20),
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    50.0,
			Description: "desc test",
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

		result, err := repo.Create(ctx, item)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when general database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()

		item := &models.SaleItem{
			SaleID:      int64(10),
			ProductID:   int64(20),
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    50.0,
			Description: "desc test",
		}

		// Simulando erro genérico do banco
		dbError := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		result, err := repo.Create(ctx, item)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})
}

func TestItemSale_Update(t *testing.T) {
	t.Run("successfully update item sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		expectedTime := time.Now()

		item := &models.SaleItem{
			ID:          int64(1),
			SaleID:      int64(10),
			ProductID:   int64(20),
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    50.0,
			Description: "desc test updated",
		}

		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				expectedTime, // updated_at
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{
				item.SaleID,
				item.ProductID,
				item.Quantity,
				item.UnitPrice,
				item.Discount,
				item.Tax,
				item.Subtotal,
				item.Description,
				item.ID,
			}).
			Return(mockRow)

		err := repo.Update(ctx, item)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, item.UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when item does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()

		item := &models.SaleItem{
			ID:          int64(999),
			SaleID:      int64(10),
			ProductID:   int64(20),
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    50.0,
			Description: "desc test",
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		err := repo.Update(ctx, item)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()

		item := &models.SaleItem{
			ID:          int64(1),
			SaleID:      int64(999), // SaleID que não existe
			ProductID:   int64(20),
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    50.0,
			Description: "desc test",
		}

		fkError := &pgconn.PgError{
			Code:    "23503",
			Message: "violação de chave estrangeira",
		}

		mockRow := &mockDb.MockRow{Err: fkError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		err := repo.Update(ctx, item)

		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when general database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()

		item := &models.SaleItem{
			ID:          int64(1),
			SaleID:      int64(10),
			ProductID:   int64(20),
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    50.0,
			Description: "desc test",
		}

		dbError := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		err := repo.Update(ctx, item)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})
}

func TestItemSale_Delete(t *testing.T) {
	t.Run("successfully delete item sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		id := int64(1)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 1, // 1 linha afetada = delete bem-sucedido
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{id}).
			Return(mockResult, nil)

		err := repo.Delete(ctx, id)

		assert.NoError(t, err)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when item does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		id := int64(999)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 0, // 0 linhas afetadas = item não existe
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{id}).
			Return(mockResult, nil)

		err := repo.Delete(ctx, id)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		id := int64(1)
		dbError := errors.New("connection lost")

		mockDB.
			On("Exec", ctx, mock.Anything, []any{id}).
			Return(nil, dbError)

		err := repo.Delete(ctx, id)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})
}

func TestItemSale_DeleteBySaleID(t *testing.T) {
	t.Run("successfully delete items by sale id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(10)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 3, // 3 itens deletados
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.DeleteBySaleID(ctx, saleID)

		assert.NoError(t, err)

		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete when no items found for sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(999)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 0, // Nenhum item encontrado para deletar
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.DeleteBySaleID(ctx, saleID)

		assert.NoError(t, err) // Não retorna erro mesmo quando não encontra itens

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(10)
		dbError := errors.New("connection lost")

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(nil, dbError)

		err := repo.DeleteBySaleID(ctx, saleID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})
}
