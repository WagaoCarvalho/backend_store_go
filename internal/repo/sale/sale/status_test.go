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

func TestSaleRepo_Cancel(t *testing.T) {
	t.Run("successfully cancel sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 1,
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.Cancel(ctx, saleID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when sale does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(999)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 0,
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.Cancel(ctx, saleID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)
		dbError := errors.New("connection failed")

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(nil, dbError)

		err := repo.Cancel(ctx, saleID)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSaleRepo_Complete(t *testing.T) {
	t.Run("successfully complete sale", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 1,
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.Complete(ctx, saleID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when sale does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(999)

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 0,
		}

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(mockResult, nil)

		err := repo.Complete(ctx, saleID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(1)
		dbError := errors.New("connection failed")

		mockDB.
			On("Exec", ctx, mock.Anything, []any{saleID}).
			Return(nil, dbError)

		err := repo.Complete(ctx, saleID)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
