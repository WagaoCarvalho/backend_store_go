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

func TestSupplierContactRelationRepo_HasSupplierContactRelation(t *testing.T) {
	t.Run("return true when relation exists", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		contactID := int64(2)

		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				1, // exists
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, contactID}).Return(mockRow)

		exists, err := repo.HasSupplierContactRelation(ctx, supplierID, contactID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return false when relation does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		contactID := int64(2)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, contactID}).Return(mockRow)

		exists, err := repo.HasSupplierContactRelation(ctx, supplierID, contactID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		contactID := int64(2)

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, contactID}).Return(mockRow)

		exists, err := repo.HasSupplierContactRelation(ctx, supplierID, contactID)

		assert.False(t, exists)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})
}
