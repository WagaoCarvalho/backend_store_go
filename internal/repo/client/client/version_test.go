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

func TestClient_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		// Criar um MockRow personalizado para int
		mockRow := &mockDb.MockRowWithInt{IntValue: 5}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 5, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("db error")
		mockRow := &mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.Zero(t, version)
		assert.ErrorIs(t, err, errMsg.ErrGetVersion)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no rows are found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		// Mock para simular QueryRow.Scan retornando pgx.ErrNoRows
		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.Zero(t, version)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

}
