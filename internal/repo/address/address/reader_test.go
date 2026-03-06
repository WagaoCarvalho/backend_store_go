package repo

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =======================
// Testes para GetByID
// =======================

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

// =======================
// Testes para GetByUserID
// =======================

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
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

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

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
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

		mockDB.On("Query", ctx, mock.Anything, []interface{}{userID}).Return(mockRows, nil)

		result, err := repo.GetByUserID(ctx, userID)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows.Err returns error after loop", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		rowsErr := errors.New("rows iteration error")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{userID}).Return(mockRows, nil)

		result, err := repo.GetByUserID(ctx, userID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("multiple addresses appended to results", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Times(3)
		mockRows.On("Scan", mock.AnythingOfType("*int64"),
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Times(3)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{userID}).Return(mockRows, nil)

		result, err := repo.GetByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

// =======================
// Testes para GetByClientCpfID
// =======================

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
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

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
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{clientCpfID}).Return(mockRows, nil)

		result, err := repo.GetByClientCpfID(ctx, clientCpfID)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		clientCpfID := int64(1)
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

		mockDB.On("Query", ctx, mock.Anything, []interface{}{clientCpfID}).Return(mockRows, nil)

		result, err := repo.GetByClientCpfID(ctx, clientCpfID)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows.Err returns error after loop", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		clientCpfID := int64(1)
		rowsErr := errors.New("rows iteration error")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{clientCpfID}).Return(mockRows, nil)

		result, err := repo.GetByClientCpfID(ctx, clientCpfID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("single address appended to empty results", func(t *testing.T) {
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
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{clientCpfID}).Return(mockRows, nil)

		result, err := repo.GetByClientCpfID(ctx, clientCpfID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

// =======================
// Testes para GetBySupplierID
// =======================

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
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return empty slice when no addresses found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
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
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when rows.Err returns error after loop", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}
		ctx := context.Background()
		supplierID := int64(1)
		rowsErr := errors.New("rows iteration error")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockRows.On("Conn").Return(nil).Maybe()
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{supplierID}).Return(mockRows, nil)

		result, err := repo.GetBySupplierID(ctx, supplierID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

// =======================
// Testes para isAllowedAddressField
// =======================

func TestIsAllowedAddressField(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected bool
	}{
		// Campos permitidos
		{"user_id", "user_id", true},
		{"client_cpf_id", "client_cpf_id", true},
		{"supplier_id", "supplier_id", true},

		// Campos NÃO permitidos (default case)
		{"id", "id", false},
		{"street", "street", false},
		{"city", "city", false},
		{"empty string", "", false},
		{"random field", "random_field", false},
		{"case sensitive USER_ID", "USER_ID", false},
		{"case sensitive User_Id", "User_Id", false},
		{"with spaces", "user id", false},
		{"sql injection attempt", "user_id; DROP TABLE", false},
		{"another sql injection", "1=1", false},
		{"OR injection", "OR 1=1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedAddressField(tt.field)
			assert.Equal(t, tt.expected, result,
				"isAllowedAddressField(%q) = %v, want %v", tt.field, result, tt.expected)
		})
	}
}

// =======================
// Testes para scanAddress
// =======================

func TestScanAddress(t *testing.T) {
	t.Run("successful scan", func(t *testing.T) {
		mockScanner := new(mockDb.MockRows)
		now := time.Now()

		mockScanner.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			if id, ok := args[0].(*int64); ok {
				*id = int64(1)
			}

			if userID, ok := args[1].(**int64); ok {
				val := int64(100)
				*userID = &val
			}

			if clientCpfID, ok := args[2].(**int64); ok {
				val := int64(200)
				*clientCpfID = &val
			}

			if supplierID, ok := args[3].(**int64); ok {
				val := int64(300)
				*supplierID = &val
			}

			if street, ok := args[4].(*string); ok {
				*street = "Rua Teste"
			}

			if streetNumber, ok := args[5].(*string); ok {
				*streetNumber = "123"
			}

			if complement, ok := args[6].(*string); ok {
				*complement = "Apto 101"
			}

			if city, ok := args[7].(*string); ok {
				*city = "São Paulo"
			}

			if state, ok := args[8].(*string); ok {
				*state = "SP"
			}

			if country, ok := args[9].(*string); ok {
				*country = "Brasil"
			}

			if postalCode, ok := args[10].(*string); ok {
				*postalCode = "01234-567"
			}

			if isActive, ok := args[11].(*bool); ok {
				*isActive = true
			}

			if createdAt, ok := args[12].(*time.Time); ok {
				*createdAt = now
			}

			if updatedAt, ok := args[13].(*time.Time); ok {
				*updatedAt = now
			}
		}).Return(nil)

		addr, err := scanAddress(mockScanner)

		assert.NoError(t, err)
		assert.NotNil(t, addr)
		assert.Equal(t, int64(1), addr.ID)
		assert.NotNil(t, addr.UserID)
		assert.Equal(t, int64(100), *addr.UserID)
		assert.NotNil(t, addr.ClientCpfID)
		assert.Equal(t, int64(200), *addr.ClientCpfID)
		assert.NotNil(t, addr.SupplierID)
		assert.Equal(t, int64(300), *addr.SupplierID)
		assert.Equal(t, "Rua Teste", addr.Street)
		assert.Equal(t, "123", addr.StreetNumber)
		assert.Equal(t, "Apto 101", addr.Complement)
		assert.Equal(t, "São Paulo", addr.City)
		assert.Equal(t, "SP", addr.State)
		assert.Equal(t, "Brasil", addr.Country)
		assert.Equal(t, "01234-567", addr.PostalCode)
		assert.True(t, addr.IsActive)
		assert.Equal(t, now, addr.CreatedAt)
		assert.Equal(t, now, addr.UpdatedAt)

		mockScanner.AssertExpectations(t)
	})

	t.Run("scan returns error", func(t *testing.T) {
		mockScanner := new(mockDb.MockRows)
		expectedErr := errors.New("scan failed")

		mockScanner.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(expectedErr)

		addr, err := scanAddress(mockScanner)

		assert.Nil(t, addr)
		assert.ErrorIs(t, err, expectedErr)
		mockScanner.AssertExpectations(t)
	})

	t.Run("scan with nil pointers", func(t *testing.T) {
		mockScanner := new(mockDb.MockRows)
		now := time.Now()

		mockScanner.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			if id, ok := args[0].(*int64); ok {
				*id = int64(2)
			}

			if userID, ok := args[1].(**int64); ok {
				*userID = nil
			}

			if clientCpfID, ok := args[2].(**int64); ok {
				*clientCpfID = nil
			}

			if supplierID, ok := args[3].(**int64); ok {
				*supplierID = nil
			}

			if street, ok := args[4].(*string); ok {
				*street = "Outra Rua"
			}

			if isActive, ok := args[11].(*bool); ok {
				*isActive = false
			}

			if createdAt, ok := args[12].(*time.Time); ok {
				*createdAt = now
			}

			if updatedAt, ok := args[13].(*time.Time); ok {
				*updatedAt = now
			}
		}).Return(nil)

		addr, err := scanAddress(mockScanner)

		assert.NoError(t, err)
		assert.NotNil(t, addr)
		assert.Equal(t, int64(2), addr.ID)
		assert.Nil(t, addr.UserID)
		assert.Nil(t, addr.ClientCpfID)
		assert.Nil(t, addr.SupplierID)
		assert.Equal(t, "Outra Rua", addr.Street)
		assert.False(t, addr.IsActive)
		assert.Equal(t, now, addr.CreatedAt)
		assert.Equal(t, now, addr.UpdatedAt)

		mockScanner.AssertExpectations(t)
	})
}

// =======================
// Testes para getByField usando reflection
// =======================

func TestGetByField_InvalidField(t *testing.T) {
	t.Run("Invalid field - should return ErrInvalidField and not query DB", func(t *testing.T) {
		t.Parallel()

		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}

		ctx := context.Background()
		field := "invalid_field"
		value := int64(123)

		result, err := repo.getByField(ctx, field, value)

		require.Nil(t, result)
		require.Error(t, err)
		require.ErrorIs(t, err, errMsg.ErrInvalidField)

		// Garantia forte: nenhuma chamada ao banco
		mockDB.AssertNotCalled(t, "Query")
		mockDB.AssertNotCalled(t, "QueryRow")
	})
	// Testar a linha if !isAllowedAddressField(field) usando reflection
	t.Run("test invalid field path using reflection", func(t *testing.T) {
		// Criar uma instância do repositório
		mockDB := new(mockDb.MockDatabase)
		repo := &addressRepo{db: mockDB}

		// Usar reflection para chamar getByField
		reflectValue := reflect.ValueOf(repo)
		method := reflectValue.MethodByName("getByField")

		if !method.IsValid() {
			t.Skip("Método getByField não encontrado (pode ser privado)")
			return
		}

		// Preparar argumentos
		ctx := context.Background()
		invalidField := "invalid_field"
		value := int64(123)

		// Chamar o método via reflection
		results := method.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(invalidField),
			reflect.ValueOf(value),
		})

		// Verificar resultados
		assert.Len(t, results, 2)

		// Primeiro resultado deve ser nil
		assert.True(t, results[0].IsNil(), "Resultado deve ser nil para campo inválido")

		// Segundo resultado deve ser erro
		err, ok := results[1].Interface().(error)
		assert.True(t, ok, "Segundo resultado deve ser error")
		assert.ErrorIs(t, err, errMsg.ErrInvalidField)

		// O mock DB não deve ser chamado quando o campo é inválido
		mockDB.AssertNotCalled(t, "Query")
	})
}

// =======================
// Teste de integração para cobertura completa
// =======================

func TestAddressRepository_CompleteCoverage(t *testing.T) {
	t.Run("test switch case coverage for isAllowedAddressField", func(t *testing.T) {
		// Testar cada caso do switch explicitamente
		testCases := []struct {
			field    string
			expected bool
			caseName string
		}{
			{"user_id", true, "case 'user_id'"},
			{"client_cpf_id", true, "case 'client_cpf_id'"},
			{"supplier_id", true, "case 'supplier_id'"},
			{"any_other", false, "default case"},
			{"id", false, "default case for 'id'"},
			{"", false, "default case for empty string"},
		}

		for _, tc := range testCases {
			t.Run(tc.caseName, func(t *testing.T) {
				result := isAllowedAddressField(tc.field)
				assert.Equal(t, tc.expected, result,
					"Field %s: expected %v, got %v", tc.field, tc.expected, result)
			})
		}
	})

	t.Run("test results append with varying counts", func(t *testing.T) {
		testCases := []struct {
			name          string
			rowCount      int
			expectedCount int
		}{
			{"no rows", 0, 0},
			{"one row", 1, 1},
			{"two rows", 2, 2},
			{"five rows", 5, 5},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockDB := new(mockDb.MockDatabase)
				repo := &addressRepo{db: mockDB}
				ctx := context.Background()
				userID := int64(1)

				mockRows := new(mockDb.MockRows)

				// Configurar Next para retornar true N vezes
				if tc.rowCount > 0 {
					mockRows.On("Next").Return(true).Times(tc.rowCount)
					mockRows.On("Scan", mock.AnythingOfType("*int64"),
						mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(nil).Times(tc.rowCount)
				}
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Err").Return(nil)
				mockRows.On("Close").Return()

				mockRows.On("Conn").Return(nil).Maybe()
				mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{}).Maybe()

				mockDB.On("Query", ctx, mock.Anything, []interface{}{userID}).Return(mockRows, nil)

				result, err := repo.GetByUserID(ctx, userID)

				assert.NoError(t, err)
				assert.Len(t, result, tc.expectedCount)

				mockDB.AssertExpectations(t)
				mockRows.AssertExpectations(t)
			})
		}
	})
}
