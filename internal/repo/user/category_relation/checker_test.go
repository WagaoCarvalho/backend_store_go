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

func TestUserCategoryRelationRepo_HasUserCategoryRelation(t *testing.T) {
	t.Run("return true when relation exists", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		categoryID := int64(2)

		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				1, // exists
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, categoryID}).Return(mockRow)

		exists, err := repo.HasUserCategoryRelation(ctx, userID, categoryID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false when relation does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		categoryID := int64(2)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, categoryID}).Return(mockRow)

		exists, err := repo.HasUserCategoryRelation(ctx, userID, categoryID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		categoryID := int64(2)

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, categoryID}).Return(mockRow)

		exists, err := repo.HasUserCategoryRelation(ctx, userID, categoryID)

		assert.False(t, exists)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrRelationCheck)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrRelationCheck.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return true when relation exists using legacy mode", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		categoryID := int64(2)

		mockRow := &mockDb.MockRow{
			Value: 1, // Usando modo legado com Value
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, categoryID}).Return(mockRow)

		exists, err := repo.HasUserCategoryRelation(ctx, userID, categoryID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)
		categoryID := int64(2)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, categoryID}).Return(mockRow)

		exists, err := repo.HasUserCategoryRelation(ctx, userID, categoryID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false with zero category ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		categoryID := int64(0)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, categoryID}).Return(mockRow)

		exists, err := repo.HasUserCategoryRelation(ctx, userID, categoryID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false with negative IDs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)
		categoryID := int64(-2)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, categoryID}).Return(mockRow)

		exists, err := repo.HasUserCategoryRelation(ctx, userID, categoryID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})
}
