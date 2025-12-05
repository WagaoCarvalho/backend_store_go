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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierRepo_Create(t *testing.T) {
	t.Run("successfully create supplier", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		supplier := &models.Supplier{
			Name:        "Test Supplier",
			CNPJ:        utils.StrToPtr("12345678000195"),
			CPF:         nil,
			Description: "Test Description",
			Status:      true,
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []interface{}{
				int64(1), // id
				now,      // created_at
				now,      // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, supplier)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		supplier := &models.Supplier{
			Name:        "Test Supplier",
			CNPJ:        utils.StrToPtr("12345678000195"),
			CPF:         nil,
			Description: "Test Description",
			Status:      true,
		}

		dbError := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, supplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierRepo_Update(t *testing.T) {
	t.Run("successfully update supplier", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		supplier := &models.Supplier{
			ID:          1,
			Name:        "Updated Supplier",
			CNPJ:        utils.StrToPtr("12345678000195"),
			CPF:         nil,
			Description: "Updated Description",
			Status:      true,
			Version:     1,
		}

		mockResult := pgconn.NewCommandTag("UPDATE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{
			supplier.Name,
			supplier.CNPJ,
			supplier.CPF,
			supplier.Description,
			supplier.Status,
			supplier.ID,
			supplier.Version,
		}).Return(mockResult, nil)

		err := repo.Update(ctx, supplier)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrVersionConflict when no rows affected", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		supplier := &models.Supplier{
			ID:          1,
			Name:        "Updated Supplier",
			CNPJ:        utils.StrToPtr("12345678000195"),
			CPF:         nil,
			Description: "Updated Description",
			Status:      true,
			Version:     1,
		}

		mockResult := pgconn.NewCommandTag("UPDATE 0")
		mockDB.On("Exec", ctx, mock.Anything, mock.Anything).Return(mockResult, nil)

		err := repo.Update(ctx, supplier)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		supplier := &models.Supplier{
			ID:          1,
			Name:        "Updated Supplier",
			CNPJ:        utils.StrToPtr("12345678000195"),
			CPF:         nil,
			Description: "Updated Description",
			Status:      true,
			Version:     1,
		}

		dbError := errors.New("database error")
		mockDB.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, dbError)

		err := repo.Update(ctx, supplier)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierRepo_Delete(t *testing.T) {
	t.Run("successfully delete supplier", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockResult := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(mockResult, nil)

		err := repo.Delete(ctx, supplierID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockResult := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(mockResult, nil)

		err := repo.Delete(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		dbError := errors.New("database error")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(nil, dbError)

		err := repo.Delete(ctx, supplierID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrDelete.Error())
		mockDB.AssertExpectations(t)
	})
}
