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

func TestClientCpfRepo_Disable(t *testing.T) {

	t.Run("successfully disable client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				clientID, // returning id
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{false, clientID}).
			Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when client does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{false, clientID}).
			Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return wrapped error on database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{false, clientID}).
			Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, dbErr)
		mockDB.AssertExpectations(t)
	})
}

func TestClientCpfRepo_Enable(t *testing.T) {

	t.Run("successfully enable client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				clientID, // returning id
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{true, clientID}).
			Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when client does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{true, clientID}).
			Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return wrapped error on database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{true, clientID}).
			Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, dbErr)
		mockDB.AssertExpectations(t)
	})
}
