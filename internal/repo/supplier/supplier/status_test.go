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

func TestSupplierRepo_Disable(t *testing.T) {
	t.Run("successfully disable supplier", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRow{
			Value: expectedTime,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, false}).Return(mockRow)

		err := repo.Disable(ctx, supplierID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, false}).Return(mockRow)

		err := repo.Disable(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDisable when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, false}).Return(mockRow)

		err := repo.Disable(ctx, supplierID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDisable)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierRepo_Enable(t *testing.T) {
	t.Run("successfully enable supplier", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRow{
			Value: expectedTime,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, true}).Return(mockRow)

		err := repo.Enable(ctx, supplierID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, true}).Return(mockRow)

		err := repo.Enable(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrEnable when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{supplierID, true}).Return(mockRow)

		err := repo.Enable(ctx, supplierID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrEnable)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}
