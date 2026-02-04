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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewContactTx(t *testing.T) {
	t.Run("deve criar nova instância de ContactTx", func(t *testing.T) {
		result := NewContactTx()

		assert.NotNil(t, result)

		_, ok := result.(*contactTx)
		assert.True(t, ok)
	})

	t.Run("deve retornar instâncias diferentes a cada chamada", func(t *testing.T) {
		instance1 := NewContactTx()
		instance2 := NewContactTx()

		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
		assert.NotSame(t, instance1, instance2)
	})
}

func TestContactTx_CreateTx(t *testing.T) {
	t.Run("sucesso ao criar contato dentro da transação", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &contactTx{}
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

		mockTx.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, contact)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)

		mockTx.AssertExpectations(t)
	})

	t.Run("deve retornar ErrNilEntity quando contact for nil", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &contactTx{}
		ctx := context.Background()

		result, err := repo.CreateTx(ctx, mockTx, nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNilModel)
	})

	t.Run("deve retornar ErrCreate em erro genérico de banco", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &contactTx{}
		ctx := context.Background()

		dbErr := errors.New("database failure")

		mockRow := &mockDb.MockRow{
			Err: dbErr,
		}

		mockTx.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, &models.Contact{})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())

		mockTx.AssertExpectations(t)
	})

	t.Run("deve retornar ErrCreate em erro pgx", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &contactTx{}
		ctx := context.Background()

		pgxErr := &pgconn.PgError{
			Message: "unique violation",
		}

		mockRow := &mockDb.MockRow{
			Err: pgxErr,
		}

		mockTx.
			On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, &models.Contact{})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, pgxErr.Message)

		mockTx.AssertExpectations(t)
	})
}
