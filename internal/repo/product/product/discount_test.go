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

func TestProductRepo_EnableDiscount(t *testing.T) {
	t.Run("successfully enable discount for product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		mockRow := &mockDb.MockRowWithInt{IntValue: 2} // Nova versão

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		err := repo.EnableDiscount(ctx, productID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		err := repo.EnableDiscount(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrProductEnableDiscount when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		err := repo.EnableDiscount(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrProductEnableDiscount)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_DisableDiscount(t *testing.T) {
	t.Run("successfully disable discount for product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		mockRow := &mockDb.MockRowWithInt{IntValue: 2} // Nova versão

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		err := repo.DisableDiscount(ctx, productID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		err := repo.DisableDiscount(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrProductDisableDiscount when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		err := repo.DisableDiscount(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrProductDisableDiscount)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_ApplyDiscount(t *testing.T) {
	t.Run("successfully apply discount to product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		discountPercent := 15.0

		// Mock para a query principal
		mockRow := &mockDb.MockRowWithID{
			IDValue:   productID,
			TimeValue: time.Now(),
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, discountPercent}).Return(mockRow)

		result, err := repo.ApplyDiscount(ctx, productID, discountPercent)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, productID, result.ID)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)
		discountPercent := 15.0

		// Mock para a query principal retornando no rows
		mockRowMain := &mockDb.MockRow{Err: pgx.ErrNoRows}
		// Mock para a query de verificação também retornando no rows
		mockRowCheck := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, discountPercent}).Return(mockRowMain)
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRowCheck)

		result, err := repo.ApplyDiscount(ctx, productID, discountPercent)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrProductDiscountNotAllowed when discount not allowed", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		discountPercent := 15.0

		// Mock para a query principal retornando no rows (allow_discount = FALSE)
		mockRowMain := &mockDb.MockRow{Err: pgx.ErrNoRows}
		// Mock para a query de verificação confirmando que produto existe
		mockRowCheck := &mockDb.MockRowWithInt{IntValue: 1}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, discountPercent}).Return(mockRowMain)
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRowCheck)

		result, err := repo.ApplyDiscount(ctx, productID, discountPercent)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrProductDiscountNotAllowed)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrProductApplyDiscount when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		discountPercent := 15.0

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, discountPercent}).Return(mockRow)

		result, err := repo.ApplyDiscount(ctx, productID, discountPercent)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrProductApplyDiscount)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
