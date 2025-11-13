package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryRelationRepo_GetAllRelationsByProductID(t *testing.T) {
	t.Run("successfully get all relations by product id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*time.Time"),
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{productID}).Return(mockRows, nil)

		result, err := repo.GetAllRelationsByProductID(ctx, productID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		dbErr := errors.New("database connection failed")
		mockDB.On("Query", ctx, mock.Anything, []interface{}{productID}).Return(nil, dbErr)

		result, err := repo.GetAllRelationsByProductID(ctx, productID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{productID}).Return(mockRows, nil)

		result, err := repo.GetAllRelationsByProductID(ctx, productID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{productID}).Return(mockRows, nil)

		result, err := repo.GetAllRelationsByProductID(ctx, productID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
