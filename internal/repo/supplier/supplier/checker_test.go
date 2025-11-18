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

func TestSupplierRepo_SupplierExists(t *testing.T) {
	t.Run("successfully check supplier exists - true", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRow := &mockDb.MockRow{
			Value: true, // Supplier exists
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		exists, err := repo.SupplierExists(ctx, supplierID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully check supplier exists - false", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockRow := &mockDb.MockRow{
			Value: false, // Supplier does not exist
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRow)

		exists, err := repo.SupplierExists(ctx, supplierID)

		assert.NoError(t, err)
		assert.False(t, exists)
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

		exists, err := repo.SupplierExists(ctx, supplierID)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "erro ao buscar")
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})
}
