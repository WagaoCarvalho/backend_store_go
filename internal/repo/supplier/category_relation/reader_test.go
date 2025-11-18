package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierCategoryRelationRepo_GetBySupplierID(t *testing.T) {
	t.Run("successfully get all relations by supplier id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		dbErr := errors.New("database connection failed")
		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(nil, dbErr)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), scanErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrScan when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), rowsErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get empty list when no relations exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get multiple relations", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRows := new(mockDb.MockRows)
		// Simula 3 relações
		mockRows.On("Next").Return(true).Times(3)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(nil).Times(3)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationRepo_GetByCategoryID(t *testing.T) {
	t.Run("successfully get all relations by category id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRows, nil)

		result, err := repo.GetByCategoryID(ctx, categoryID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		dbErr := errors.New("database connection failed")
		mockDB.On("Query", ctx, mock.Anything, []interface{}{categoryID}).Return(nil, dbErr)

		result, err := repo.GetByCategoryID(ctx, categoryID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRows, nil)

		result, err := repo.GetByCategoryID(ctx, categoryID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), scanErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrScan when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRows, nil)

		result, err := repo.GetByCategoryID(ctx, categoryID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), rowsErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get empty list when no relations exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRows, nil)

		result, err := repo.GetByCategoryID(ctx, categoryID)

		assert.NoError(t, err)
		assert.NotNil(t, result) // Agora será uma slice vazia, não nil
		assert.Empty(t, result)
		assert.Len(t, result, 0)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get multiple relations", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		mockRows := new(mockDb.MockRows)
		// Simula 3 relações
		mockRows.On("Next").Return(true).Times(3)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(nil).Times(3)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRows, nil)

		result, err := repo.GetByCategoryID(ctx, categoryID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
