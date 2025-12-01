package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContact_Create(t *testing.T) {
	t.Run("successfully create contact", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()

		cont := &models.Contact{
			ContactName:        "Maria Souza",
			ContactDescription: "Gerente de Vendas",
			Email:              "maria@empresa.com",
			Phone:              "1133224455",
			Cell:               "11987654321",
			ContactType:        "Comercial",
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithID{
			IDValue:   1,
			TimeValue: now,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, cont)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate when unique constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{}

		pgErr := &pgconn.PgError{
			Code:    "23505",
			Message: "duplicate key value violates unique constraint",
		}
		mockRow := &mockDb.MockRow{
			Err: pgErr,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, cont)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})
	t.Run("return ErrCreate when check constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{}

		pgErr := &pgconn.PgError{
			Code:    "23514",
			Message: "check constraint violation",
		}

		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, cont)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, pgErr.Message)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{}

		pgErr := &pgconn.PgError{
			Code:    "23503",
			Message: "foreign key violation",
		}

		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, cont)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})

}

func TestContact_Update(t *testing.T) {
	t.Run("successfully update contact", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()

		cont := &models.Contact{
			ID:                 1,
			ContactName:        "Maria Souza",
			ContactDescription: "Gerente de Vendas",
			Email:              "maria@empresa.com",
			Phone:              "1133224455",
			Cell:               "11987654321",
			ContactType:        "Comercial",
		}

		now := time.Now()
		mockRow := &mockDb.MockRow{Value: now}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, cont)

		assert.NoError(t, err)
		assert.Equal(t, now, cont.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no rows updated", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{ID: 999}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, cont)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate when unique constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{ID: 1}

		pgErr := &pgconn.PgError{
			Code:    "23505",
			Message: "duplicate key value violates unique constraint",
		}
		mockRow := &mockDb.MockRow{Err: pgErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, cont)

		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData when check constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{ID: 1}

		pgErr := &pgconn.PgError{
			Code:    "23514",
			Message: "check constraint violation",
		}
		mockRow := &mockDb.MockRow{Err: pgErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, cont)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when pg error different from known cases", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{ID: 1}

		pgErr := &pgconn.PgError{
			Code:    "23503",
			Message: "foreign key violation",
		}
		mockRow := &mockDb.MockRow{Err: pgErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, cont)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, pgErr.Message)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when query fails with generic error", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		cont := &models.Contact{ID: 1}

		dbErr := errors.New("timeout")
		mockRow := &mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, cont)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestContact_Delete(t *testing.T) {
	t.Run("successfully delete contact", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		contactID := int64(1)

		mockRow := &mockDb.MockRowWithID{IDValue: contactID}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{contactID}).
			Return(mockRow)

		err := repo.Delete(ctx, contactID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when contact not found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		contactID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{contactID}).
			Return(mockRow)

		err := repo.Delete(ctx, contactID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		contactID := int64(1)

		dbErr := errors.New("connection lost")
		mockRow := &mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{contactID}).
			Return(mockRow)

		err := repo.Delete(ctx, contactID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}
