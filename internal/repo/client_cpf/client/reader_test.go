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

func TestClientCpfRepo_GetByID(t *testing.T) {

	t.Run("successfully get client by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		clientID := int64(1)
		now := time.Now()

		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				clientID,          // id
				"Cliente Teste",   // name
				"teste@email.com", // email
				"12345678909",     // cpf
				"descrição",       // description
				true,              // status
				1,                 // version
				now,               // created_at
				now,               // updated_at
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).
			Return(mockRow)

		result, err := repo.GetByID(ctx, clientID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, clientID, result.ID)
		assert.Equal(t, "Cliente Teste", result.Name)
		assert.Equal(t, "teste@email.com", result.Email)
		assert.Equal(t, "12345678909", result.CPF)
		assert.Equal(t, 1, result.Version)
		assert.True(t, result.Status)
		assert.Equal(t, now, result.CreatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when client does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		clientID := int64(999)

		mockRow := &mockDb.MockRow{
			Err: pgx.ErrNoRows,
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).
			Return(mockRow)

		result, err := repo.GetByID(ctx, clientID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)

		mockDB.AssertExpectations(t)
	})

	t.Run("return wrapped error on database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		clientID := int64(1)
		dbErr := errors.New("database failure")

		mockRow := &mockDb.MockRow{
			Err: dbErr,
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).
			Return(mockRow)

		result, err := repo.GetByID(ctx, clientID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbErr)

		mockDB.AssertExpectations(t)
	})
}
