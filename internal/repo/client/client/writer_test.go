package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_Create(t *testing.T) {
	t.Run("successfully create client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			Name:        "Test Client",
			Email:       utils.StrToPtr("test@example.com"),
			CPF:         utils.StrToPtr("123.456.789-00"),
			CNPJ:        utils.StrToPtr("12.345.678/0001-90"),
			Description: "Test description",
			Status:      true,
			Version:     1,
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithID{
			IDValue:   1,
			TimeValue: now,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate when unique constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData when check constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		pgErr := &pgconn.PgError{Code: "23514", Message: "check constraint violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		dbErr := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other pg error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, pgErr.Message)
		mockDB.AssertExpectations(t)
	})
}

func TestClient_Update(t *testing.T) {
	t.Run("successfully update client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:          1,
			Name:        "Updated Client",
			Email:       utils.StrToPtr("updated@example.com"),
			CPF:         utils.StrToPtr("987.654.321-00"),
			CNPJ:        utils.StrToPtr("98.765.432/0001-10"),
			Status:      true,
			Description: "Updated description",
			Version:     1,
		}

		mockRow := &mockDb.MockRowWithInt{IntValue: 2}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, client)

		assert.NoError(t, err)
		assert.Equal(t, 2, client.Version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrVersionConflict when no rows affected", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{ID: 1, Version: 1}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate when unique constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{ID: 1, Version: 1}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData when check constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{ID: 1, Version: 1}

		pgErr := &pgconn.PgError{Code: "23514", Message: "check constraint violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when other pg error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{ID: 1, Version: 1}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, pgErr.Message)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when generic database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{ID: 1, Version: 1}

		dbErr := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("successfully delete client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		cmdTag := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{clientID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		dbError := errors.New("database connection failed")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{clientID}).Return(pgconn.CommandTag{}, dbError)

		err := repo.Delete(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
