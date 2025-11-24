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

func TestUserRepo_UserExists(t *testing.T) {
	t.Run("successfully check user exists - true", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		mockRow := &mockDb.MockRow{
			Value: true, // User exists
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		exists, err := repo.UserExists(ctx, userID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully check user exists - false", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockRow := &mockDb.MockRow{
			Value: false, // User does not exist
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		exists, err := repo.UserExists(ctx, userID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		exists, err := repo.UserExists(ctx, userID)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "erro ao buscar")
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully check user exists with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)

		mockRow := &mockDb.MockRow{
			Value: false, // User with ID 0 does not exist
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		exists, err := repo.UserExists(ctx, userID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully check user exists with negative user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)

		mockRow := &mockDb.MockRow{
			Value: false, // User with negative ID does not exist
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		exists, err := repo.UserExists(ctx, userID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})
}
