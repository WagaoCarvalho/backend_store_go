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

func TestUserRepo_GetByID(t *testing.T) {
	t.Run("successfully get user by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(1),           // UID (id na query)
				"testuser",         // username
				"test@example.com", // email
				"Test description", // description
				true,               // status
				createdAt,          // created_at
				updatedAt,          // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(1), user.UID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test description", user.Description)
		assert.Equal(t, true, user.Status)
		assert.Equal(t, createdAt, user.CreatedAt)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, user)
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

		user, err := repo.GetByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "erro ao buscar")
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get user with false status", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(2)

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(2),                    // UID
				"inactiveuser",              // username
				"inactive@example.com",      // email
				"Inactive user description", // description
				false,                       // status
				createdAt,                   // created_at
				updatedAt,                   // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(2), user.UID)
		assert.Equal(t, "inactiveuser", user.Username)
		assert.Equal(t, "inactive@example.com", user.Email)
		assert.Equal(t, "Inactive user description", user.Description)
		assert.Equal(t, false, user.Status)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get user with empty description", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(3)

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(3),            // UID
				"user3",             // username
				"user3@example.com", // email
				"",                  // empty description
				true,                // status
				createdAt,           // created_at
				updatedAt,           // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(3), user.UID)
		assert.Equal(t, "user3", user.Username)
		assert.Equal(t, "user3@example.com", user.Email)
		assert.Equal(t, "", user.Description)
		assert.Equal(t, true, user.Status)
		mockDB.AssertExpectations(t)
	})
}
