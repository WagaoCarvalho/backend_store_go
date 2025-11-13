package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryRelationRepo_HasProductCategoryRelation(t *testing.T) {
	t.Run("return true when relation exists", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		categoryID := int64(2)

		mockRow := &mockDb.MockRowWithInt{IntValue: 1}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, categoryID}).Return(mockRow)

		exists, err := repo.HasProductCategoryRelation(ctx, productID, categoryID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false when relation does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		categoryID := int64(2)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, categoryID}).Return(mockRow)

		exists, err := repo.HasProductCategoryRelation(ctx, productID, categoryID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)
		categoryID := int64(2)

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{productID, categoryID}).Return(mockRow)

		exists, err := repo.HasProductCategoryRelation(ctx, productID, categoryID)

		assert.False(t, exists)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
