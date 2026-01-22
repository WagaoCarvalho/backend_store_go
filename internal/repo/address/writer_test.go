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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddress_Create(t *testing.T) {
	t.Run("successfully create address", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()

		addr := &models.Address{
			UserID:       utils.Int64Ptr(1),
			ClientCpfID:  utils.Int64Ptr(2),
			SupplierID:   utils.Int64Ptr(3),
			Street:       "Rua A",
			StreetNumber: "123",
			City:         "São Paulo",
			State:        "SP",
			Country:      "Brasil",
			PostalCode:   "01000-000",
			IsActive:     true,
		}

		now := time.Now()
		mockRow := mockDb.MockRow{Value: now}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow, nil)

		result, err := repo.Create(ctx, addr)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey on FK violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addr := &models.Address{}

		fkErr := errMsgPg.NewForeignKeyViolation("addresses_user_id_fkey")

		mockRow := mockDb.MockRow{Err: fkErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow, nil)

		result, err := repo.Create(ctx, addr)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addr := &models.Address{}
		dbErr := errors.New("database failure")

		mockRow := mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow, nil)

		result, err := repo.Create(ctx, addr)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestAddress_Update(t *testing.T) {
	t.Run("successfully update address", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		now := time.Now()

		addr := &models.Address{
			ID:           1,
			UserID:       utils.Int64Ptr(1),
			ClientCpfID:  utils.Int64Ptr(2),
			SupplierID:   utils.Int64Ptr(3),
			Street:       "Rua Nova",
			StreetNumber: "456",
			City:         "Curitiba",
			State:        "PR",
			Country:      "Brasil",
			PostalCode:   "80000-000",
			IsActive:     true,
		}

		mockRow := mockDb.MockRow{Value: now}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow, nil)

		err := repo.Update(ctx, addr)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when address doesn't exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addr := &models.Address{}

		mockRow := mockDb.MockRow{Err: errors.New("no rows in result set")}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow, nil)

		err := repo.Update(ctx, addr)

		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no rows found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()

		address := &models.Address{
			ID:           999,
			UserID:       utils.Int64Ptr(1),
			ClientCpfID:  utils.Int64Ptr(2),
			SupplierID:   utils.Int64Ptr(3),
			Street:       "Rua Teste",
			StreetNumber: "123",
			City:         "São Paulo",
			State:        "SP",
			Country:      "Brasil",
			PostalCode:   "01000-000",
			IsActive:     true,
		}

		// Simula QueryRow retornando ErrNoRows
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockDb.MockRow{Err: pgx.ErrNoRows}, nil)

		err := repo.Update(ctx, address)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when FK violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addr := &models.Address{}

		fkErr := errMsgPg.NewForeignKeyViolation("addresses_client_cpf_id_fkey")
		mockRow := mockDb.MockRow{Err: fkErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow, nil)

		err := repo.Update(ctx, addr)

		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate on database error", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addr := &models.Address{}
		dbErr := errors.New("timeout")

		mockRow := mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow, nil)

		err := repo.Update(ctx, addr)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestAddress_Delete(t *testing.T) {
	t.Run("successfully delete address", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()

		cmdTag := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(cmdTag, nil)

		err := repo.Delete(ctx, 1)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no rows deleted", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()

		cmdTag := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(cmdTag, nil)

		err := repo.Delete(ctx, 999)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete on database error", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		dbErr := errors.New("connection lost")

		mockDB.On("Exec", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(pgconn.CommandTag{}, dbErr)

		err := repo.Delete(ctx, 1)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}
