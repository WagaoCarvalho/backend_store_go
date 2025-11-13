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

func TestProductRepo_GetByID(t *testing.T) {
	t.Run("successfully get product by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		mockRow := &mockDb.MockRow{}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		product, err := repo.GetByID(ctx, productID)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		product, err := repo.GetByID(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, product)
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

		product, err := repo.GetByID(ctx, productID)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		// MockRow com valor específico
		mockRow := &mockDb.MockRow{
			Value: int64(5), // Valor que será retornado no Scan
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, productID)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Zero(t, version)
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

		version, err := repo.GetVersionByID(ctx, productID)

		assert.Error(t, err)
		assert.Zero(t, version)
		assert.Contains(t, err.Error(), "erro ao buscar versão")
		mockDB.AssertExpectations(t)
	})
}
