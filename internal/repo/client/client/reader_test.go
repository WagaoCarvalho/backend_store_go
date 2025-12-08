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

func TestClient_GetByID(t *testing.T) {
	t.Run("successfully get client by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRowWithID{
			IDValue:   clientID,
			TimeValue: expectedTime,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).Return(mockRow)

		result, err := repo.GetByID(ctx, clientID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTime, result.CreatedAt)
		assert.Equal(t, clientID, result.ID)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when client does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).Return(mockRow)

		result, err := repo.GetByID(ctx, clientID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)
		dbError := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).Return(mockRow)

		result, err := repo.GetByID(ctx, clientID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetByName(t *testing.T) {
	t.Run("successfully get clients by name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("*string"),    // Name
			mock.AnythingOfType("**string"),   // Email (ponteiro para ponteiro)
			mock.AnythingOfType("**string"),   // Phone (ponteiro para ponteiro)
			mock.AnythingOfType("**string"),   // Document (ponteiro para ponteiro)
			mock.AnythingOfType("*string"),    // Address
			mock.AnythingOfType("*bool"),      // IsActive
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, "test")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("db failure")

		// Retornar um MockRows vazio junto com o erro, ou apenas o erro
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(&mockDb.MockRows{}, dbErr)

		result, err := repo.GetByName(ctx, "test")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, "test")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
