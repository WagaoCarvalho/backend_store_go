package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	address "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// Testes do método Filter
// ============================================================================

func TestAddressFilter_Filter_Success_WithDifferentCombinations(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	// Configura rows vazias
	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	userID := int64(1)
	clientCpfID := int64(2)
	supplierID := int64(3)

	filters := []*address.Address{
		{IsActive: true},
		{UserID: &userID, IsActive: true},
		{ClientCpfID: &clientCpfID, IsActive: true},
		{SupplierID: &supplierID, IsActive: true},
		{City: "São Paulo", IsActive: true},
		{State: "sp", IsActive: true},
		{PostalCode: "01234567", IsActive: true},
		{PostalCode: "01234-567", IsActive: true},
		{PostalCode: "01.234-567", IsActive: true},
		{Street: "Rua Teste", IsActive: true},
		{StreetNumber: "123", IsActive: true},
		{Country: "Brasil", IsActive: true},
	}

	for _, f := range filters {
		result, err := repo.Filter(ctx, f)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	}
}

func TestAddressFilter_Filter_WithResults(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	now := time.Now()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Run(func(args mock.Arguments) {
		// Preenche todos os campos
		if id, ok := args[0].(*int64); ok {
			*id = 1
		}
		if userID, ok := args[1].(*int64); ok {
			*userID = 100
		}
		if clientCpfID, ok := args[2].(*int64); ok {
			*clientCpfID = 200
		}
		if supplierID, ok := args[3].(*int64); ok {
			*supplierID = 300
		}
		if street, ok := args[4].(*string); ok {
			*street = "Rua Teste"
		}
		if streetNumber, ok := args[5].(*string); ok {
			*streetNumber = "123"
		}
		if complement, ok := args[6].(**string); ok {
			comp := "Apto 101"
			*complement = &comp
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
			*postalCode = "01234567"
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
	}).Return(nil).Once()
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, &address.Address{IsActive: true})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "Rua Teste", result[0].Street)
	assert.Equal(t, "Apto 101", result[0].Complement)
	assert.Equal(t, "São Paulo", result[0].City)
}

func TestAddressFilter_Filter_WithNilComplement(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Run(func(args mock.Arguments) {
		if complement, ok := args[6].(**string); ok {
			*complement = nil
		}
	}).Return(nil).Once()
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, &address.Address{IsActive: true})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Empty(t, result[0].Complement)
}

func TestAddressFilter_Filter_QueryError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	dbErr := errors.New("database connection error")
	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

	result, err := repo.Filter(ctx, &address.Address{IsActive: true})

	assert.Error(t, err)
	assert.ErrorIs(t, err, errMsg.ErrGet)
	assert.Contains(t, err.Error(), "database connection error")
	assert.Nil(t, result)
}

func TestAddressFilter_Filter_ScanError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Return(errors.New("scan error"))
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, &address.Address{IsActive: true})

	assert.Error(t, err)
	assert.ErrorIs(t, err, errMsg.ErrScan)
	assert.Contains(t, err.Error(), "scan error")
	assert.Nil(t, result)
}

func TestAddressFilter_Filter_RowsError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(errors.New("rows iteration error"))
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, &address.Address{IsActive: true})

	assert.Error(t, err)
	assert.ErrorIs(t, err, errMsg.ErrIterate)
	assert.Contains(t, err.Error(), "rows iteration error")
	assert.Nil(t, result)
}

// ============================================================================
// Testes do método FindActive
// ============================================================================

func TestAddressFilter_FindActive_Success(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	filter := &address.Address{City: "São Paulo"}
	result, err := repo.FindActive(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.True(t, filter.IsActive)
}

func TestAddressFilter_FindActive_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result, err := repo.FindActive(ctx, &address.Address{})

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================================
// Testes do método FindByPostalCode
// ============================================================================

func TestAddressFilter_FindByPostalCode_Success(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	testCases := []struct {
		code  string
		exact bool
	}{
		{"01234567", true},
		{"01234567", false},
		{"01234-567", true},
		{"01234-567", false},
		{"01.234-567", true},
		{"01.234-567", false},
		{"", true},
		{"", false},
	}

	for _, tc := range testCases {
		result, err := repo.FindByPostalCode(ctx, tc.code, tc.exact)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	}
}

func TestAddressFilter_FindByPostalCode_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result, err := repo.FindByPostalCode(ctx, "01234567", true)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Teste específico para cobrir o trecho "if result == nil" no FindByPostalCode
func TestAddressFilter_FindByPostalCode_NilResult(t *testing.T) {
	// Similar ao FindActive, o Filter nunca retorna nil devido ao make
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.FindByPostalCode(ctx, "01234567", true)

	assert.NoError(t, err)
	assert.NotNil(t, result) // Não é nil, então o if não é executado
	assert.Empty(t, result)
}

// ============================================================================
// Testes do método FindByCityAndState
// ============================================================================

func TestAddressFilter_FindByCityAndState_Success(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	testCases := []struct {
		city  string
		state string
	}{
		{"São Paulo", "SP"},
		{"Rio de Janeiro", "RJ"},
		{"", ""},
		{" ", " "},
	}

	for _, tc := range testCases {
		result, err := repo.FindByCityAndState(ctx, tc.city, tc.state)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	}
}

func TestAddressFilter_FindByCityAndState_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result, err := repo.FindByCityAndState(ctx, "São Paulo", "SP")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Teste específico para cobrir o trecho "if result == nil" no FindByCityAndState
func TestAddressFilter_FindByCityAndState_NilResult(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.FindByCityAndState(ctx, "São Paulo", "SP")

	assert.NoError(t, err)
	assert.NotNil(t, result) // Não é nil, então o if não é executado
	assert.Empty(t, result)
}

// ============================================================================
// Testes do método FindByPostalCodeV2 (versão com função auxiliar)
// ============================================================================

func TestAddressFilter_FindByPostalCodeV2_Success(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.FindByPostalCodeV2(ctx, "01234567", true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestAddressFilter_FindByPostalCodeV2_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result, err := repo.FindByPostalCodeV2(ctx, "01234567", true)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Teste específico para cobrir o retorno de ensureNonNilSlice no FindByPostalCodeV2
func TestAddressFilter_FindByPostalCodeV2_WithEnsureNonNilSlice(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	// O Filter retorna slice vazia, ensureNonNilSlice retorna a mesma slice
	result, err := repo.FindByPostalCodeV2(ctx, "01234567", true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// ============================================================================
// Testes da função auxiliar ensureNonNilSlice
// ============================================================================

func TestEnsureNonNilSlice(t *testing.T) {
	addr1 := &address.Address{ID: 1}
	addr2 := &address.Address{ID: 2}

	tests := []struct {
		name  string
		input []*address.Address
	}{
		{"slice nil", nil},
		{"slice vazia", []*address.Address{}},
		{"um elemento", []*address.Address{addr1}},
		{"dois elementos", []*address.Address{addr1, addr2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureNonNilSlice(tt.input)

			assert.NotNil(t, result)
			assert.Len(t, result, len(tt.input))

			for i := range tt.input {
				assert.Equal(t, tt.input[i], result[i])
			}
		})
	}
}

// ============================================================================
// Teste para o mapa de campos de ordenação
// ============================================================================

func TestAddressAllowedSortFields(t *testing.T) {
	expectedFields := []string{
		"id", "user_id", "client_cpf_id", "supplier_id",
		"street", "street_number", "city", "state",
		"country", "postal_code", "is_active", "created_at", "updated_at",
	}

	for _, field := range expectedFields {
		t.Run(field, func(t *testing.T) {
			mappedField, exists := addressAllowedSortFields[field]
			assert.True(t, exists, "Campo %s deve existir no mapa", field)
			assert.Equal(t, field, mappedField)
		})
	}

	invalidFields := []string{"invalid", "name", "email", ""}
	for _, field := range invalidFields {
		t.Run("invalid_"+field, func(t *testing.T) {
			_, exists := addressAllowedSortFields[field]
			assert.False(t, exists, "Campo %s não deve existir no mapa", field)
		})
	}
}
