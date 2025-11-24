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

func TestUserRepo_Disable(t *testing.T) {
	t.Run("successfully disable user", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				2,         // version (int)
				updatedAt, // updated_at (time.Time)
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Disable(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Disable(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Disable(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDisable)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrDisable.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully disable user with zero ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				1,         // version
				updatedAt, // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Disable(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with negative user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Disable(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})
}

func TestUserRepo_Enable(t *testing.T) {
	t.Run("successfully enable user", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				3,         // version (int)
				updatedAt, // updated_at (time.Time)
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Enable(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Enable(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Enable(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrEnable)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrEnable.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully enable user with zero ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				1,         // version
				updatedAt, // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Enable(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with negative user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Enable(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully enable user with high version number", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(5)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				100,       // version (número alto)
				updatedAt, // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		err := repo.Enable(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
