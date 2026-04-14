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

func TestSupplier_GetByID(t *testing.T) {
	t.Run("successfully get supplier by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRow := &mockDb.MockRow{}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		result, err := repo.GetByID(ctx, supplierID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		result, err := repo.GetByID(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		result, err := repo.GetByID(ctx, supplierID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		mockDB.AssertExpectations(t)
	})
}

func TestSupplier_GetByName(t *testing.T) {
	t.Run("successfully get suppliers by name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		name := "Fornecedor"

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = 1
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "Fornecedor A"
			}
			if ptr, ok := args.Get(2).(**string); ok {
				cnpj := "12345678000195"
				*ptr = &cnpj
			}
			if ptr, ok := args.Get(5).(*int); ok {
				*ptr = 1
			}
			if ptr, ok := args.Get(6).(*bool); ok {
				*ptr = true
			}
			if ptr, ok := args.Get(7).(*time.Time); ok {
				*ptr = time.Now()
			}
			if ptr, ok := args.Get(8).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, name)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no suppliers found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		name := "nonexistent"

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, name)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		dbErr := errors.New("database error")
		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(nil, dbErr)

		result, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		scanErr := errors.New("scan error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		rowsErr := errors.New("rows error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSupplier_GetAll(t *testing.T) {
	t.Run("successfully get all suppliers", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		// Primeiro supplier
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = 1
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "Supplier 1"
			}
			if ptr, ok := args.Get(2).(**string); ok {
				cnpj := "12345678000195"
				*ptr = &cnpj
			}
			if ptr, ok := args.Get(5).(*int); ok {
				*ptr = 1
			}
			if ptr, ok := args.Get(6).(*bool); ok {
				*ptr = true
			}
			if ptr, ok := args.Get(7).(*time.Time); ok {
				*ptr = time.Now()
			}
			if ptr, ok := args.Get(8).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()

		// Segundo supplier
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = 2
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "Supplier 2"
			}
			if ptr, ok := args.Get(3).(**string); ok {
				cpf := "12345678901"
				*ptr = &cpf
			}
			if ptr, ok := args.Get(5).(*int); ok {
				*ptr = 2
			}
			if ptr, ok := args.Get(6).(*bool); ok {
				*ptr = false
			}
			if ptr, ok := args.Get(7).(*time.Time); ok {
				*ptr = time.Now().Add(-24 * time.Hour)
			}
			if ptr, ok := args.Get(8).(*time.Time); ok {
				*ptr = time.Now().Add(-24 * time.Hour)
			}
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Supplier 1", result[0].Name)
		assert.Equal(t, int64(2), result[1].ID)
		assert.Equal(t, "Supplier 2", result[1].Name)
		assert.Equal(t, false, result[1].Status)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return empty slice when no suppliers found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("database error")
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

		result, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSupplier_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		expectedVersion := int64(5)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		result, err := repo.GetVersionByID(ctx, supplierID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		result, err := repo.GetVersionByID(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		result, err := repo.GetVersionByID(ctx, supplierID)

		assert.Error(t, err)
		assert.Equal(t, int64(0), result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
	})
}
