package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierRepo_Disable(t *testing.T) {
	t.Run("successfully disable supplier", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockResult := pgconn.NewCommandTag("UPDATE 1")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(mockResult, nil)

		err := repo.Disable(ctx, supplierID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockResult := pgconn.NewCommandTag("UPDATE 0")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(mockResult, nil)

		err := repo.Disable(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database exec fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		execErr := errors.New("database error")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(nil, execErr)

		err := repo.Disable(ctx, supplierID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDisable)
		assert.Contains(t, err.Error(), execErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrDisable.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierRepo_Enable(t *testing.T) {
	t.Run("successfully enable supplier", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockResult := pgconn.NewCommandTag("UPDATE 1")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(mockResult, nil)

		err := repo.Enable(ctx, supplierID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when supplier does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(999)

		mockResult := pgconn.NewCommandTag("UPDATE 0")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(mockResult, nil)

		err := repo.Enable(ctx, supplierID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database exec fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		execErr := errors.New("database error")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{supplierID}).Return(nil, execErr)

		err := repo.Enable(ctx, supplierID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrEnable)
		assert.Contains(t, err.Error(), execErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrEnable.Error())
		mockDB.AssertExpectations(t)
	})
}
