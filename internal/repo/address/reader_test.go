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

func TestAddress_GetByID(t *testing.T) {
	t.Run("successfully get address by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(1)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRowWithID{
			IDValue:   addressID,
			TimeValue: expectedTime,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{addressID}).Return(mockRow)

		result, err := repo.GetByID(ctx, addressID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTime, result.CreatedAt)
		assert.Equal(t, addressID, result.ID)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when address does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{addressID}).Return(mockRow)

		result, err := repo.GetByID(ctx, addressID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		addressID := int64(1)
		dbError := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{addressID}).Return(mockRow)

		result, err := repo.GetByID(ctx, addressID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestAddress_GetByUserID(t *testing.T) {
	t.Run("successfully get addresses by user id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"),
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return()

		// ✅ Tornar opcionais
		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{userID}).Return(mockRows, nil)

		result, err := repo.GetByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		dbError := errors.New("db down")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{userID}).Return(mockRows, dbError)

		result, err := repo.GetByUserID(ctx, userID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})

}

func TestAddress_GetByClientCpfID(t *testing.T) {
	t.Run("successfully get addresses by client_cpf id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		clientCpfID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"),
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return()

		// ✅ Tornar opcionais
		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{clientCpfID}).Return(mockRows, nil)

		result, err := repo.GetByClientCpfID(ctx, clientCpfID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return empty slice when no addresses found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		clientCpfID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Close").Return()

		// ✅ Tornar opcionais
		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{clientCpfID}).Return(mockRows, nil)

		result, err := repo.GetByClientCpfID(ctx, clientCpfID)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestAddress_GetBySupplierID(t *testing.T) {
	t.Run("successfully get addresses by supplier id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Twice()
		mockRows.On("Scan", mock.AnythingOfType("*int64"),
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return()

		// ✅ Tornar opcionais
		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		scanErr := errors.New("scan failed")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true)
		mockRows.On("Scan", mock.AnythingOfType("*int64"),
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(scanErr)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, scanErr.Error()) // ✅ apenas isso
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

}
