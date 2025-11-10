package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewContactTx(t *testing.T) {
	t.Run("successfully create new ContactTx instance", func(t *testing.T) {
		var db *pgxpool.Pool

		result := NewContactTx(db)

		assert.NotNil(t, result)
		_, ok := result.(*contactTx)
		assert.True(t, ok, "Expected result to be of type *contactTx")
	})

	t.Run("return instance with provided db pool", func(t *testing.T) {
		var db *pgxpool.Pool
		result := NewContactTx(db)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		var db *pgxpool.Pool

		instance1 := NewContactTx(db)
		instance2 := NewContactTx(db)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestContactTx_CreateTx(t *testing.T) {
	t.Run("successfully create contact within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &contactTx{db: &pgxpool.Pool{}}
		ctx := context.Background()

		contact := &models.Contact{
			ContactName:        "João Silva",
			ContactDescription: "Gerente de compras",
			Email:              "joao.silva@empresa.com",
			Phone:              "1133224455",
			Cell:               "11999887766",
			ContactType:        "Comercial",
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithID{
			IDValue:   1,
			TimeValue: now,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, contact)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when query fails with generic database error", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &contactTx{db: &pgxpool.Pool{}}
		ctx := context.Background()
		contact := &models.Contact{}
		dbErr := errors.New("database connection failed")

		mockRow := &mockDb.MockRow{
			Err: dbErr,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, contact)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when pgx error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &contactTx{db: &pgxpool.Pool{}}
		ctx := context.Background()
		contact := &models.Contact{}
		pgxErr := &pgconn.PgError{
			Message: "unique violation",
		}

		mockRow := &mockDb.MockRow{
			Err: pgxErr,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, contact)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, pgxErr.Message)
		mockTx.AssertExpectations(t)
	})
}
