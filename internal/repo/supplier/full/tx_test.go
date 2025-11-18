package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewSupplierFull(t *testing.T) {
	t.Run("successfully create new SupplierFull instance", func(t *testing.T) {
		var db repo.DBTransactor

		result := NewSupplierFull(db)

		assert.NotNil(t, result)

		_, ok := result.(*supplierFullRepo)
		assert.True(t, ok, "Expected result to be of type *supplierFullRepo")
	})

	t.Run("return instance with provided db transactor", func(t *testing.T) {
		var db repo.DBTransactor

		result := NewSupplierFull(db)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		var db repo.DBTransactor

		instance1 := NewSupplierFull(db)
		instance2 := NewSupplierFull(db)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestSupplierFullRepo_BeginTx(t *testing.T) {
	t.Run("successfully begin transaction", func(t *testing.T) {
		mockDB := new(mockDb.MockDBTransactor)
		repo := &supplierFullRepo{db: mockDB}
		ctx := context.Background()

		mockTx := new(mockDb.MockTx)
		mockDB.On("BeginTx", ctx, pgx.TxOptions{}).Return(mockTx, nil)

		tx, err := repo.BeginTx(ctx)

		assert.NoError(t, err)
		assert.Equal(t, mockTx, tx)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when begin transaction fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDBTransactor)
		repo := &supplierFullRepo{db: mockDB}
		ctx := context.Background()

		dbError := errors.New("transaction failed")
		// Retornar MockTx vazio junto com o erro
		mockDB.On("BeginTx", ctx, pgx.TxOptions{}).Return(&mockDb.MockTx{}, dbError)

		tx, err := repo.BeginTx(ctx)

		assert.NotNil(t, tx) // Agora retorna um MockTx (mesmo que vazio)
		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierFullRepo_CreateTx(t *testing.T) {
	t.Run("successfully create supplier within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierFullRepo{db: nil} // db não é usado no CreateTx
		ctx := context.Background()

		supplier := &models.Supplier{
			Name:   "Test Supplier",
			CNPJ:   utils.StrToPtr("12345678000195"),
			CPF:    utils.StrToPtr(""),
			Status: true,
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []interface{}{
				int64(1),  // id
				1,         // version
				createdAt, // created_at
				updatedAt, // updated_at
			},
		}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{
			supplier.Name,
			supplier.CNPJ,
			supplier.CPF,
			supplier.Status,
		}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, supplier)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, supplier, result)
		assert.Equal(t, int64(1), supplier.ID)
		assert.Equal(t, 1, supplier.Version)
		assert.Equal(t, createdAt, supplier.CreatedAt)
		assert.Equal(t, updatedAt, supplier.UpdatedAt)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierFullRepo{db: nil}
		ctx := context.Background()

		supplier := &models.Supplier{
			Name:   "Test Supplier",
			CNPJ:   utils.StrToPtr("12345678000195"),
			CPF:    utils.StrToPtr(""),
			Status: true,
		}

		dbError := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbError}
		mockTx.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, supplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockTx.AssertExpectations(t)
	})
}
