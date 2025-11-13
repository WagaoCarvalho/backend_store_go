package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductRepo_ProductExists(t *testing.T) {
	t.Run("successfully check product exists - true", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		mockRow := &mockDb.MockRow{
			Value: true, // Product exists
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		exists, err := repo.ProductExists(ctx, productID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully check product exists - false", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		mockRow := &mockDb.MockRow{
			Value: false, // Product does not exist
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID}).Return(mockRow)

		exists, err := repo.ProductExists(ctx, productID)

		assert.NoError(t, err)
		assert.False(t, exists)
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

		exists, err := repo.ProductExists(ctx, productID)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockDB.AssertExpectations(t)
	})
}
