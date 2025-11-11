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

func TestContact_GetByID(t *testing.T) {
	t.Run("successfully get contact by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		contactID := int64(1)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRowWithID{
			IDValue:   contactID,
			TimeValue: expectedTime,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{contactID}).Return(mockRow)

		result, err := repo.GetByID(ctx, contactID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTime, result.CreatedAt)
		assert.Equal(t, contactID, result.ID)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when contact does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		contactID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{contactID}).Return(mockRow)

		result, err := repo.GetByID(ctx, contactID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &contactRepo{db: mockDB}
		ctx := context.Background()
		contactID := int64(1)
		dbError := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{contactID}).Return(mockRow)

		result, err := repo.GetByID(ctx, contactID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
