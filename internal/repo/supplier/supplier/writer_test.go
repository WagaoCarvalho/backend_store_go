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
	"github.com/jackc/pgx/v5"
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
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(1), // id
				now,      // created_at
				now,      // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 5 &&
				args[0] == supplier.Name &&
				args[1] == supplier.CNPJ &&
				args[2] == supplier.CPF &&
				args[3] == supplier.Description &&
				args[4] == supplier.Status
		})).Return(mockRow)

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

		mockDB.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(mockRow)

		result, err := repo.Create(ctx, supplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), dbError.Error())
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

		now := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				now, // updated_at
				2,   // new version
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 7 &&
				args[0] == supplier.Name &&
				args[1] == supplier.CNPJ &&
				args[2] == supplier.CPF &&
				args[3] == supplier.Description &&
				args[4] == supplier.Status &&
				args[5] == supplier.ID &&
				args[6] == supplier.Version
		})).Return(mockRow)

		err := repo.Update(ctx, supplier)

		assert.NoError(t, err)
		assert.Equal(t, now, supplier.UpdatedAt)
		assert.Equal(t, 2, supplier.Version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrVersionConflict when version mismatch", func(t *testing.T) {
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

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		// Mock para verificar se o registro existe
		existsRow := &mockDb.MockRow{
			Values: []interface{}{true},
		}
		mockDB.On("QueryRow", ctx, "SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1)", []interface{}{supplier.ID}).Return(existsRow)

		mockDB.On("QueryRow", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 7 && args[5] == supplier.ID && args[6] == supplier.Version
		})).Return(mockRow)

		err := repo.Update(ctx, supplier)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		supplier := &models.Supplier{
			ID:          999,
			Name:        "Updated Supplier",
			CNPJ:        utils.StrToPtr("12345678000195"),
			CPF:         nil,
			Description: "Updated Description",
			Status:      true,
			Version:     1,
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		// Mock para verificar que o registro NÃO existe
		existsRow := &mockDb.MockRow{
			Values: []interface{}{false},
		}
		mockDB.On("QueryRow", ctx, "SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1)", []interface{}{supplier.ID}).Return(existsRow)

		mockDB.On("QueryRow", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 7 && args[5] == supplier.ID
		})).Return(mockRow)

		err := repo.Update(ctx, supplier)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when check exists query fails", func(t *testing.T) {
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

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		// Mock para verificar se o registro existe - RETORNA ERRO
		checkError := errors.New("check exists database error")
		existsRow := &mockDb.MockRow{Err: checkError}
		mockDB.On("QueryRow", ctx, "SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1)", []interface{}{supplier.ID}).Return(existsRow)

		mockDB.On("QueryRow", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 7 && args[5] == supplier.ID && args[6] == supplier.Version
		})).Return(mockRow)

		err := repo.Update(ctx, supplier)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), checkError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when check exists query fails with connection error", func(t *testing.T) {
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

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		// Mock para verificar se o registro existe - RETORNA ERRO DE CONEXÃO
		connectionError := errors.New("connection refused")
		existsRow := &mockDb.MockRow{Err: connectionError}
		mockDB.On("QueryRow", ctx, "SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1)", []interface{}{supplier.ID}).Return(existsRow)

		mockDB.On("QueryRow", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 7 && args[5] == supplier.ID && args[6] == supplier.Version
		})).Return(mockRow)

		err := repo.Update(ctx, supplier)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), connectionError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when update query fails with database error", func(t *testing.T) {
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
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(mockRow)

		err := repo.Update(ctx, supplier)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), dbError.Error())
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
		mockDB.AssertExpectations(t)
	})
}
