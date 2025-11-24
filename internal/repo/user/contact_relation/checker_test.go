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

func TestUserContactRelationRepo_HasUserContactRelation(t *testing.T) {
	t.Run("return true when relation exists", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userContactRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		contactID := int64(2)

		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				1, // exists
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, contactID}).Return(mockRow)

		exists, err := repo.HasUserContactRelation(ctx, userID, contactID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false when relation does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userContactRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		contactID := int64(2)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, contactID}).Return(mockRow)

		exists, err := repo.HasUserContactRelation(ctx, userID, contactID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userContactRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		contactID := int64(2)

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, contactID}).Return(mockRow)

		exists, err := repo.HasUserContactRelation(ctx, userID, contactID)

		assert.False(t, exists)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return true when relation exists using legacy mode", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userContactRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		contactID := int64(2)

		mockRow := &mockDb.MockRow{
			Value: 1, // Usando modo legado com Value
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, contactID}).Return(mockRow)

		exists, err := repo.HasUserContactRelation(ctx, userID, contactID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false when relation does not exist with custom error message", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userContactRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		contactID := int64(2)

		// Testando com a string exata que a função verifica
		mockRow := &mockDb.MockRow{Err: errors.New("no rows in result set")}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, contactID}).Return(mockRow)

		exists, err := repo.HasUserContactRelation(ctx, userID, contactID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userContactRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)
		contactID := int64(2)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, contactID}).Return(mockRow)

		exists, err := repo.HasUserContactRelation(ctx, userID, contactID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error with zero contact ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userContactRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		contactID := int64(0)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID, contactID}).Return(mockRow)

		exists, err := repo.HasUserContactRelation(ctx, userID, contactID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})
}
