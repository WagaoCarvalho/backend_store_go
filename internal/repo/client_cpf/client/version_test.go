package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetVersionByID(t *testing.T) {
	const query = `
		SELECT version
		FROM clients_cpf
		WHERE id = $1
		LIMIT 1
	`

	t.Run("successfully get version by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		mockRow := &mockDb.MockRowWithInt{IntValue: 5}
		mockDB.On("QueryRow", ctx, query, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 5, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return raw error when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("db error")
		mockRow := &mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, query, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.Zero(t, version)
		assert.ErrorIs(t, err, dbErr)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no rows are found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, query, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.Zero(t, version)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})
}
