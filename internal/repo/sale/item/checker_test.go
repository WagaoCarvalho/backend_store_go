package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItemSale_ItemExists(t *testing.T) {
	t.Run("return true when item exists", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		id := int64(1)

		mockRow := &mockDb.MockRow{
			Value: true, // EXISTS retorna true
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{id}).
			Return(mockRow)

		exists, err := repo.ItemExists(ctx, id)

		assert.NoError(t, err)
		assert.True(t, exists)

		mockDB.AssertExpectations(t)
	})

	t.Run("return false when item does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		id := int64(999)

		mockRow := &mockDb.MockRow{
			Value: false, // EXISTS retorna false
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{id}).
			Return(mockRow)

		exists, err := repo.ItemExists(ctx, id)

		assert.NoError(t, err)
		assert.False(t, exists)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		id := int64(1)
		dbError := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{id}).
			Return(mockRow)

		exists, err := repo.ItemExists(ctx, id)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})
}
