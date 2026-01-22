package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientCredit_GetByID(t *testing.T) {
	t.Run("successfully get client credit by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		creditID := int64(1)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRowWithID{
			IDValue:   creditID,
			TimeValue: expectedTime,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{creditID}).Return(mockRow)

		result, err := repo.GetByID(ctx, creditID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTime, result.CreatedAt)
		assert.Equal(t, creditID, result.ID)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when client credit not found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		creditID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{creditID}).Return(mockRow)

		result, err := repo.GetByID(ctx, creditID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database connection fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		creditID := int64(2)
		dbErr := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{creditID}).Return(mockRow)

		result, err := repo.GetByID(ctx, creditID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when Scan returns pgconn.PgError", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		creditID := int64(3)
		pgErr := &pgconn.PgError{Message: "syntax error"}

		mockRow := &mockDb.MockRow{Err: pgErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{creditID}).Return(mockRow)

		result, err := repo.GetByID(ctx, creditID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, "syntax error")
		mockDB.AssertExpectations(t)
	})
}

func TestClientCredit_GetByClientID(t *testing.T) {
	t.Run("successfully get client credit by client id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(10)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRowWithID{
			IDValue:   clientID,
			TimeValue: expectedTime,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).Return(mockRow)

		result, err := repo.GetByClientID(ctx, clientID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, clientID, result.ClientID)
		assert.Equal(t, expectedTime, result.CreatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when client credit not found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).Return(mockRow)

		result, err := repo.GetByClientID(ctx, clientID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database connection fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(2)
		dbErr := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).Return(mockRow)

		result, err := repo.GetByClientID(ctx, clientID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when Scan returns pgconn.PgError", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(3)
		pgErr := &pgconn.PgError{Message: "syntax error"}

		mockRow := &mockDb.MockRow{Err: pgErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{clientID}).Return(mockRow)

		result, err := repo.GetByClientID(ctx, clientID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, "syntax error")
		mockDB.AssertExpectations(t)
	})
}

func TestClientCredit_GetAll(t *testing.T) {
	t.Run("successfully get all client credits", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"), // id
			mock.AnythingOfType("*int64"), // client_id
			mock.Anything,                 // allow_credit
			mock.Anything,                 // credit_limit
			mock.Anything,                 // credit_balance
			mock.Anything,                 // created_at
			mock.Anything,                 // updated_at
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// A função Query recebe (ctx, query string, args []interface{})
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("db fail")
		mockRows := new(mockDb.MockRows)

		// A função Query recebe (ctx, query string, args []interface{})
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, dbErr)

		result, err := repo.GetAll(ctx)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*int64"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		// A função Query recebe (ctx, query string, args []interface{})
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.GetAll(ctx)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestClientCredit_GetByName(t *testing.T) {
	t.Run("successfully get client credits by name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		name := "John"

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"), // id
			mock.AnythingOfType("*int64"), // client_id
			mock.Anything,                 // allow_credit
			mock.Anything,                 // credit_limit
			mock.Anything,                 // credit_balance
			mock.Anything,                 // created_at
			mock.Anything,                 // updated_at
		).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// Query recebe (ctx, query, args []interface{})
		mockDB.On("Query", ctx, mock.Anything, []interface{}{name}).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, name)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		name := "Maria"

		dbErr := errors.New("query failed")
		mockRows := new(mockDb.MockRows)

		mockDB.On("Query", ctx, mock.Anything, []interface{}{name}).Return(mockRows, dbErr)

		result, err := repo.GetByName(ctx, name)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		name := "Lucas"

		scanErr := errors.New("scan failed")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*int64"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{name}).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, name)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		name := "Ana"

		rowsErr := errors.New("rows iteration error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{name}).Return(mockRows, nil)

		result, err := repo.GetByName(ctx, name)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestClientCredit_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		// MockRow personalizado que retorna um int fixo
		mockRow := &mockDb.MockRowWithInt{IntValue: 3}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 3, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("db error")
		mockRow := &mockDb.MockRow{Err: dbErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{int64(1)}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, 1)

		assert.Zero(t, version)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}
