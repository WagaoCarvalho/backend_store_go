package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductRepo_GetStock(t *testing.T) {
	t.Run("successfully get stock by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		mockRow := &mockDb.MockRow{
			Value: 50, // Stock quantity
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		stock, err := repo.GetStock(ctx, productID)

		assert.NoError(t, err)
		assert.Equal(t, 50, stock)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		stock, err := repo.GetStock(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Zero(t, stock)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		stock, err := repo.GetStock(ctx, productID)

		assert.Error(t, err)
		assert.Zero(t, stock)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_UpdateStock(t *testing.T) {
	t.Run("successfully update stock", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		quantity := 100

		mockRow := &mockDb.MockRow{
			Value: 2, // New version
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, quantity}).Return(mockRow)

		err := repo.UpdateStock(ctx, productID, quantity)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error for negative quantity", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		quantity := -10

		err := repo.UpdateStock(ctx, productID, quantity)

		assert.ErrorIs(t, err, errMsg.ErrInvalidQuantity)
		mockDB.AssertNotCalled(t, "QueryRow")
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)
		quantity := 50

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, quantity}).Return(mockRow)

		err := repo.UpdateStock(ctx, productID, quantity)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		quantity := 75

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, quantity}).Return(mockRow)

		err := repo.UpdateStock(ctx, productID, quantity)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar")
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_IncreaseStock(t *testing.T) {
	t.Run("successfully increase stock", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		amount := 25

		mockRow := &mockDb.MockRow{
			Value: 3, // New version
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, amount}).Return(mockRow)

		err := repo.IncreaseStock(ctx, productID, amount)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error for zero or negative amount", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		// Teste com 0
		err := repo.IncreaseStock(ctx, productID, 0)
		assert.ErrorIs(t, err, errMsg.ErrInvalidQuantity)

		// Teste com negativo
		err = repo.IncreaseStock(ctx, productID, -5)
		assert.ErrorIs(t, err, errMsg.ErrInvalidQuantity)

		mockDB.AssertNotCalled(t, "QueryRow")
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)
		amount := 10

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, amount}).Return(mockRow)

		err := repo.IncreaseStock(ctx, productID, amount)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		amount := 15

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, amount}).Return(mockRow)

		err := repo.IncreaseStock(ctx, productID, amount)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar")
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_DecreaseStock(t *testing.T) {
	t.Run("successfully decrease stock", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		amount := 10

		mockRow := &mockDb.MockRow{
			Value: 4, // New version
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, amount}).Return(mockRow)

		err := repo.DecreaseStock(ctx, productID, amount)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error for zero or negative amount", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		// Teste com 0
		err := repo.DecreaseStock(ctx, productID, 0)
		assert.ErrorIs(t, err, errMsg.ErrInvalidQuantity)

		// Teste com negativo
		err = repo.DecreaseStock(ctx, productID, -5)
		assert.ErrorIs(t, err, errMsg.ErrInvalidQuantity)

		mockDB.AssertNotCalled(t, "QueryRow")
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)
		amount := 5

		// Mock para query principal
		mockRowMain := &mockDb.MockRow{Err: pgx.ErrNoRows}
		// Mock para verificação de existência
		mockRowCheck := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, amount}).Return(mockRowMain)
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRowCheck)

		err := repo.DecreaseStock(ctx, productID, amount)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInsufficientStock when not enough stock", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		amount := 100 // Mais do que tem no estoque

		// Mock para query principal (falha)
		mockRowMain := &mockDb.MockRow{Err: pgx.ErrNoRows}
		// Mock para verificação (produto existe)
		mockRowCheck := &mockDb.MockRow{
			Value: 1,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, amount}).Return(mockRowMain)
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRowCheck)

		err := repo.DecreaseStock(ctx, productID, amount)

		assert.ErrorIs(t, err, errMsg.ErrInsufficientStock)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		amount := 8

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, amount}).Return(mockRow)

		err := repo.DecreaseStock(ctx, productID, amount)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar")
		mockDB.AssertExpectations(t)
	})
}
