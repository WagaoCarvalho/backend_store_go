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

func TestAddress_Disable(t *testing.T) {
	t.Run("successfully disable address", func(t *testing.T) {

		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(1)
		expectedUpdatedAt := time.Now()

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockDb.MockRow{Value: expectedUpdatedAt}, nil)

		err := repo.Disable(ctx, addressID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when address does not exist", func(t *testing.T) {

		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(999)

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockDb.MockRow{Err: pgx.ErrNoRows}, nil)

		err := repo.Disable(ctx, addressID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDisable when database error occurs", func(t *testing.T) {

		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(1)
		dbError := errors.New("database connection failed")

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockDb.MockRow{Err: dbError}, nil)

		err := repo.Disable(ctx, addressID)

		assert.ErrorIs(t, err, errMsg.ErrDisable)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestAddress_Enable(t *testing.T) {
	t.Run("successfully enable address", func(t *testing.T) {

		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(1)
		expectedUpdatedAt := time.Now()

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockDb.MockRow{Value: expectedUpdatedAt}, nil)

		err := repo.Enable(ctx, addressID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when address does not exist", func(t *testing.T) {

		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(999)

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockDb.MockRow{Err: pgx.ErrNoRows}, nil)

		err := repo.Enable(ctx, addressID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrEnable when database error occurs", func(t *testing.T) {

		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(1)
		dbError := errors.New("database connection failed")

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockDb.MockRow{Err: dbError}, nil)

		err := repo.Enable(ctx, addressID)

		assert.ErrorIs(t, err, errMsg.ErrEnable)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
