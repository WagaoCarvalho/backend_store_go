package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryRelationRepo_GetAllRelationsByUserID(t *testing.T) {
	t.Run("successfully get all relations by user id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		mockRows := new(mockDb.MockRows)
		// Configurar TODOS os métodos que serão chamados
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.UserID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Run(func(args mock.Arguments) {
			// Simular o preenchimento dos valores
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(1)
			}
			if ptr, ok := args.Get(1).(*int64); ok {
				*ptr = int64(100)
			}
			if ptr, ok := args.Get(2).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once()
		mockRows.On("Close").Return().Once()

		// CORREÇÃO: Usar mock.Anything para os argumentos varargs
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAllRelationsByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].UserID)
		assert.Equal(t, int64(100), result[0].CategoryID)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return empty slice when user id is zero", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)

		result, err := repo.GetAllRelationsByUserID(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		mockDB.AssertExpectations(t)
	})

	t.Run("return empty slice when user id is negative", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)

		result, err := repo.GetAllRelationsByUserID(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		dbErr := errors.New("database connection failed")
		// CORREÇÃO: Usar mock.Anything para os argumentos varargs
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

		result, err := repo.GetAllRelationsByUserID(ctx, userID)

		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.UserID
			mock.AnythingOfType("*int64"),     // &rel.CategoryID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(scanErr).Once()
		mockRows.On("Close").Return().Once()

		// CORREÇÃO: Usar mock.Anything para os argumentos varargs
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAllRelationsByUserID(ctx, userID)

		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), scanErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr).Once()
		mockRows.On("Close").Return().Once()

		// CORREÇÃO: Usar mock.Anything para os argumentos varargs
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetAllRelationsByUserID(ctx, userID)

		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.Contains(t, err.Error(), rowsErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrIterate.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

}
