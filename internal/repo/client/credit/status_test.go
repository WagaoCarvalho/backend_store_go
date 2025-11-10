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

func TestClientCredit_Disable(t *testing.T) {
	t.Run("successfully disable client credit", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCredit{db: mockDB}
		ctx := context.Background()
		creditID := int64(1)

		cmdTag := mockDb.MockCommandTag{RowsAffectedCount: 1}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{creditID}).Return(cmdTag, nil)

		err := repo.Disable(ctx, creditID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCredit{db: mockDB}
		ctx := context.Background()
		creditID := int64(1)
		dbError := errors.New("database connection failed")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{creditID}).Return(nil, dbError)

		err := repo.Disable(ctx, creditID)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestClientCredit_Enable(t *testing.T) {
	t.Run("successfully enable client credit", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCredit{db: mockDB}
		ctx := context.Background()
		creditID := int64(1)

		cmdTag := mockDb.MockCommandTag{RowsAffectedCount: 1}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{creditID}).Return(cmdTag, nil)

		err := repo.Enable(ctx, creditID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCredit{db: mockDB}
		ctx := context.Background()
		creditID := int64(1)
		dbError := errors.New("database connection failed")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{creditID}).Return(nil, dbError)

		err := repo.Enable(ctx, creditID)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
