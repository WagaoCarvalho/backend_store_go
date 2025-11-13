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

func TestProductCategoryRepo_GetAll(t *testing.T) {
	t.Run("successfully get all product categories", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("database connection failed")
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

		result, err := repo.GetAll(ctx)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get multiple categories", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		// Simula 3 categorias
		mockRows.On("Next").Return(true).Times(3)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil).Times(3)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestProductCategoryRepo_GetByID(t *testing.T) {
	t.Run("successfully get product category by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		mockRow := &mockDb.MockRow{}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRow)

		result, err := repo.GetByID(ctx, categoryID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when category does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRow)

		result, err := repo.GetByID(ctx, categoryID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{categoryID}).Return(mockRow)

		result, err := repo.GetByID(ctx, categoryID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
