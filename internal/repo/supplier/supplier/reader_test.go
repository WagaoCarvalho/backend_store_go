package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierRepo_GetByID(t *testing.T) {
	t.Run("successfully get supplier by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRow := &mockDb.MockRow{}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		supplier, err := repo.GetByID(ctx, supplierID)

		assert.NoError(t, err)
		assert.NotNil(t, supplier)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		supplier, err := repo.GetByID(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, supplier)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		supplier, err := repo.GetByID(ctx, supplierID)

		assert.Error(t, err)
		assert.Nil(t, supplier)
		assert.Contains(t, err.Error(), "erro ao buscar")
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierRepo_GetByName(t *testing.T) {
	t.Run("successfully get suppliers by name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &s.ID
			mock.AnythingOfType("*string"),    // &s.Name
			mock.AnythingOfType("**string"),   // &s.CNPJ (ponteiro para string)
			mock.AnythingOfType("**string"),   // &s.CPF (ponteiro para string)
			mock.AnythingOfType("*string"),    // &s.Description
			mock.AnythingOfType("*bool"),      // &s.Status (bool, não string)
			mock.AnythingOfType("*time.Time"), // &s.CreatedAt
			mock.AnythingOfType("*time.Time"), // &s.UpdatedAt
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		suppliers, err := repo.GetByName(ctx, name)

		assert.NoError(t, err)
		assert.NotNil(t, suppliers)
		assert.Len(t, suppliers, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	// Os outros testes permanecem iguais, apenas corrigindo o Scan onde necessário
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

		suppliers, err := repo.GetByName(ctx, name)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, suppliers)
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

		suppliers, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, suppliers)
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
			mock.AnythingOfType("*int64"),     // &s.ID
			mock.AnythingOfType("*string"),    // &s.Name
			mock.AnythingOfType("**string"),   // &s.CNPJ
			mock.AnythingOfType("**string"),   // &s.CPF
			mock.AnythingOfType("*string"),    // &s.Description
			mock.AnythingOfType("*bool"),      // &s.Status
			mock.AnythingOfType("*time.Time"), // &s.CreatedAt
			mock.AnythingOfType("*time.Time"), // &s.UpdatedAt
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		suppliers, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, suppliers)
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

		suppliers, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, suppliers)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSupplierRepo_GetAll(t *testing.T) {
	t.Run("successfully get all suppliers", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Twice() // Dois suppliers
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &s.ID
			mock.AnythingOfType("*string"),    // &s.Name
			mock.AnythingOfType("**string"),   // &s.CNPJ
			mock.AnythingOfType("**string"),   // &s.CPF
			mock.AnythingOfType("*string"),    // &s.Description
			mock.AnythingOfType("*bool"),      // &s.Status
			mock.AnythingOfType("*time.Time"), // &s.CreatedAt
			mock.AnythingOfType("*time.Time"), // &s.UpdatedAt
		).Return(nil).Twice()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// Corrigir: usar mock.Anything para os args, pois a função pode passar nil
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		suppliers, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, suppliers)
		assert.Len(t, suppliers, 2)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	// Corrigir os outros testes também
	t.Run("return ErrNotFound when no suppliers found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		suppliers, err := repo.GetAll(ctx)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, suppliers)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("database error")
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

		suppliers, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, suppliers)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &s.ID
			mock.AnythingOfType("*string"),    // &s.Name
			mock.AnythingOfType("**string"),   // &s.CNPJ
			mock.AnythingOfType("**string"),   // &s.CPF
			mock.AnythingOfType("*string"),    // &s.Description
			mock.AnythingOfType("*bool"),      // &s.Status
			mock.AnythingOfType("*time.Time"), // &s.CreatedAt
			mock.AnythingOfType("*time.Time"), // &s.UpdatedAt
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		suppliers, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, suppliers)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()

		rowsErr := errors.New("rows error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		suppliers, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, suppliers)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSupplierRepo_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		expectedVersion := int64(5)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, supplierID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, supplierID)

		assert.Error(t, err)
		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
	})
}
