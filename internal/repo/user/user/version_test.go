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

func TestUserRepo_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		expectedVersion := int64(5)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
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

		version, err := repo.GetVersionByID(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrGetVersion)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get version zero", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(2)
		expectedVersion := int64(0)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get version with high number", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(3)
		expectedVersion := int64(100)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with negative user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
		mockDB.AssertExpectations(t)
	})
}
