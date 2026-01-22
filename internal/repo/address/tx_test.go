package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAddressTx(t *testing.T) {
	t.Run("successfully create new AddressTx instance", func(t *testing.T) {

		var db *pgxpool.Pool

		result := NewAddressTx(db)

		assert.NotNil(t, result)

		_, ok := result.(*addressTx)
		assert.True(t, ok, "Expected result to be of type *addressTx")

	})

	t.Run("return instance with provided db pool", func(t *testing.T) {

		var db *pgxpool.Pool

		result := NewAddressTx(db)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		var db *pgxpool.Pool

		instance1 := NewAddressTx(db)
		instance2 := NewAddressTx(db)

		assert.NotSame(t, instance1, instance2)

		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestAddressTx_CreateTx(t *testing.T) {
	t.Run("successfully create address within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &addressTx{db: &pgxpool.Pool{}}
		ctx := context.Background()

		address := &models.Address{
			UserID:       utils.Int64Ptr(1),
			ClientCpfID:  utils.Int64Ptr(2),
			SupplierID:   utils.Int64Ptr(3),
			Street:       "Rua A",
			StreetNumber: "123",
			Complement:   "Apto 1",
			City:         "São Paulo",
			State:        "SP",
			Country:      "Brasil",
			PostalCode:   "01000-000",
			IsActive:     true,
		}

		now := time.Now()

		mockRow := &mockDb.MockRowWithID{
			IDValue:   1,
			TimeValue: now,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, address)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey on FK violation", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &addressTx{db: &pgxpool.Pool{}}
		ctx := context.Background()
		address := &models.Address{}

		fkErr := errMsgPg.NewForeignKeyViolation("addresses_user_id_fkey")
		mockRow := &mockDb.MockRow{
			Err: fkErr,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, address)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when query fails with generic database error", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &addressTx{db: &pgxpool.Pool{}}
		ctx := context.Background()
		address := &models.Address{}
		dbErr := errors.New("database connection failed")

		mockRow := &mockDb.MockRow{
			Err: dbErr,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, address)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when pgx error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &addressTx{db: &pgxpool.Pool{}}
		ctx := context.Background()
		address := &models.Address{}
		pgxErr := &pgconn.PgError{
			Message: "connection timeout",
		}

		mockRow := &mockDb.MockRow{
			Err: pgxErr,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, address)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, pgxErr.Message)
		mockTx.AssertExpectations(t)
	})
}
