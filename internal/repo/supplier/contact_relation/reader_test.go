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

func TestSupplierContactRelationRepo_GetAllRelationsBySupplierID(t *testing.T) {
	t.Run("successfully get all relations by supplier id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRows := new(mockDb.MockRows)
		// Configurar TODOS os métodos que serão chamados
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.ContactID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Run(func(args mock.Arguments) {
			// Simular o preenchimento dos valores
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(1)
			}
			if ptr, ok := args.Get(1).(*int64); ok {
				*ptr = int64(2)
			}
			if ptr, ok := args.Get(2).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once()
		mockRows.On("Close").Return().Once()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetAllRelationsBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		dbErr := errors.New("database connection failed")
		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(nil, dbErr)

		result, err := repo.GetAllRelationsBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &rel.SupplierID
			mock.AnythingOfType("*int64"),     // &rel.ContactID
			mock.AnythingOfType("*time.Time"), // &rel.CreatedAt
		).Return(scanErr).Once()
		mockRows.On("Close").Return().Once()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetAllRelationsBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), scanErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr).Once()
		mockRows.On("Close").Return().Once()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetAllRelationsBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), rowsErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

}
