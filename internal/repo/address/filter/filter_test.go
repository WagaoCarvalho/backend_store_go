package repo

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	filterAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address/filter"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// Testes do método Filter
// ============================================================================

func TestAddressFilter_Filter_Success_WithResults(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	now := time.Now()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Run(func(args mock.Arguments) {
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

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
		IsActive:   &isActive,
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "Rua Teste", result[0].Street)
	assert.Equal(t, "123", result[0].StreetNumber)
	assert.Equal(t, "Apto 101", result[0].Complement)
	assert.Equal(t, "São Paulo", result[0].City)
	assert.Equal(t, "SP", result[0].State)
	assert.Equal(t, "Brasil", result[0].Country)
	assert.Equal(t, "01234567", result[0].PostalCode)
	assert.True(t, result[0].IsActive)
	mockDB.AssertExpectations(t)
	rows.AssertExpectations(t)
}

func TestAddressFilter_Filter_WithAllFilters(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	userID := int64(1)
	clientCpfID := int64(2)
	supplierID := int64(3)
	isActive := true
	createdFrom := time.Now().Add(-24 * time.Hour)
	createdTo := time.Now()
	updatedFrom := time.Now().Add(-12 * time.Hour)
	updatedTo := time.Now()

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "city",
			SortOrder: "asc",
			Limit:     10,
			Offset:    0,
		},
		UserID:       &userID,
		ClientCpfID:  &clientCpfID,
		SupplierID:   &supplierID,
		Street:       "Rua Teste",
		StreetNumber: "123",
		Complement:   "Apto 101",
		City:         "São Paulo",
		State:        "SP",
		Country:      "Brasil",
		PostalCode:   "01234567",
		IsActive:     &isActive,
		CreatedFrom:  &createdFrom,
		CreatedTo:    &createdTo,
		UpdatedFrom:  &updatedFrom,
		UpdatedTo:    &updatedTo,
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockDB.AssertExpectations(t)
}

func TestAddressFilter_Filter_EmptyResult(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
		IsActive:   &isActive,
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockDB.AssertExpectations(t)
}

func TestAddressFilter_Filter_MultipleRows(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)

	// Primeira linha
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Run(func(args mock.Arguments) {
		if id, ok := args[0].(*int64); ok {
			*id = 1
		}
		if street, ok := args[4].(*string); ok {
			*street = "Rua 1"
		}
		if complement, ok := args[6].(**string); ok {
			*complement = nil
		}
		if isActive, ok := args[11].(*bool); ok {
			*isActive = true
		}
	}).Return(nil).Once()

	// Segunda linha
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Run(func(args mock.Arguments) {
		if id, ok := args[0].(*int64); ok {
			*id = 2
		}
		if street, ok := args[4].(*string); ok {
			*street = "Rua 2"
		}
		if complement, ok := args[6].(**string); ok {
			*complement = nil
		}
		if isActive, ok := args[11].(*bool); ok {
			*isActive = true
		}
	}).Return(nil).Once()

	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
		IsActive:   &isActive,
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Rua 1", result[0].Street)
	assert.Equal(t, "Rua 2", result[1].Street)
	mockDB.AssertExpectations(t)
}

// ============================================================================
// Testes de ordenação
// ============================================================================

func TestAddressFilter_SortField_Valid(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "city",
			SortOrder: "asc",
			Limit:     10,
		},
		IsActive: &isActive,
	}

	mockDB.
		On("Query",
			ctx,
			mock.MatchedBy(func(q string) bool {
				return strings.Contains(q, "ORDER BY city asc")
			}),
			mock.Anything,
		).
		Return(rows, nil)

	_, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
}

func TestAddressFilter_SortField_Invalid(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "invalid_field",
			SortOrder: "asc",
			Limit:     10,
		},
		IsActive: &isActive,
	}

	mockDB.
		On("Query",
			ctx,
			mock.MatchedBy(func(q string) bool {
				return strings.Contains(q, "ORDER BY created_at asc")
			}),
			mock.Anything,
		).
		Return(rows, nil)

	_, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
}

func TestAddressFilter_SortOrder_Invalid(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "created_at",
			SortOrder: "INVALID",
			Limit:     10,
		},
		IsActive: &isActive,
	}

	mockDB.
		On("Query",
			ctx,
			mock.MatchedBy(func(q string) bool {
				return strings.Contains(q, "ORDER BY created_at desc")
			}),
			mock.Anything,
		).
		Return(rows, nil)

	_, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
}

func TestAddressFilter_SortOrder_Descending(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "city",
			SortOrder: "desc",
			Limit:     10,
		},
		IsActive: &isActive,
	}

	mockDB.
		On("Query",
			ctx,
			mock.MatchedBy(func(q string) bool {
				return strings.Contains(q, "ORDER BY city desc")
			}),
			mock.Anything,
		).
		Return(rows, nil)

	_, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
}

func TestAddressFilter_AllSortFields(t *testing.T) {
	for sortField := range addressAllowedSortFields {
		t.Run(sortField, func(t *testing.T) {
			mockDB := new(mockDb.MockDatabase)
			repo := &addressFilterRepo{db: mockDB}
			ctx := context.Background()
			isActive := true

			rows := new(mockDb.MockRows)
			rows.On("Next").Return(false)
			rows.On("Err").Return(nil)
			rows.On("Close").Return()

			f := &filterAddress.AddressFilter{
				BaseFilter: filter.BaseFilter{
					SortBy:    sortField,
					SortOrder: "asc",
					Limit:     10,
				},
				IsActive: &isActive,
			}

			mockDB.
				On("Query",
					ctx,
					mock.MatchedBy(func(q string) bool {
						expectedField := addressAllowedSortFields[sortField]
						return strings.Contains(q, "ORDER BY "+expectedField+" asc")
					}),
					mock.Anything,
				).
				Return(rows, nil)

			_, err := repo.Filter(ctx, f)
			assert.NoError(t, err)
		})
	}
}

func TestAddressFilter_CaseInsensitive(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	t.Run("sort field uppercase", func(t *testing.T) {
		f := &filterAddress.AddressFilter{
			BaseFilter: filter.BaseFilter{
				SortBy:    "CREATED_AT",
				SortOrder: "asc",
				Limit:     10,
			},
			IsActive: &isActive,
		}

		mockDB.
			On("Query",
				ctx,
				mock.MatchedBy(func(q string) bool {
					return strings.Contains(q, "ORDER BY created_at asc")
				}),
				mock.Anything,
			).
			Return(rows, nil)

		_, err := repo.Filter(ctx, f)
		assert.NoError(t, err)
	})

	t.Run("sort order uppercase", func(t *testing.T) {
		f := &filterAddress.AddressFilter{
			BaseFilter: filter.BaseFilter{
				SortBy:    "city",
				SortOrder: "DESC",
				Limit:     10,
			},
			IsActive: &isActive,
		}

		mockDB.
			On("Query",
				ctx,
				mock.MatchedBy(func(q string) bool {
					return strings.Contains(q, "ORDER BY city desc")
				}),
				mock.Anything,
			).
			Return(rows, nil)

		_, err := repo.Filter(ctx, f)
		assert.NoError(t, err)
	})
}

// ============================================================================
// Testes de erro
// ============================================================================

func TestAddressFilter_QueryError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	dbErr := errors.New("db error")
	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
		IsActive:   &isActive,
	}

	result, err := repo.Filter(ctx, f)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errMsg.ErrGet)
	assert.Contains(t, err.Error(), "db error")
}

func TestAddressFilter_ScanError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

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

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
		IsActive:   &isActive,
	}

	result, err := repo.Filter(ctx, f)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errMsg.ErrScan)
	assert.Contains(t, err.Error(), "scan error")
}

func TestAddressFilter_RowsError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()
	isActive := true

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(errors.New("iterate error"))
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	f := &filterAddress.AddressFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
		IsActive:   &isActive,
	}

	result, err := repo.Filter(ctx, f)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errMsg.ErrIterate)
	assert.Contains(t, err.Error(), "iterate error")
}

// ============================================================================
// Testes de métodos auxiliares
// ============================================================================

func TestAddressFilter_FindActive(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	filter := &filterAddress.AddressFilter{
		City: "São Paulo",
	}

	result, err := repo.FindActive(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.NotNil(t, filter.IsActive)
	assert.True(t, *filter.IsActive)
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

	t.Run("exact match true", func(t *testing.T) {
		result, err := repo.FindByPostalCode(ctx, "01234567", true)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("exact match false", func(t *testing.T) {
		result, err := repo.FindByPostalCode(ctx, "01234", false)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("CEP com formatação", func(t *testing.T) {
		result, err := repo.FindByPostalCode(ctx, "01234-567", true)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}

func TestAddressFilter_FindByPostalCode_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	dbErr := errors.New("database connection error")
	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

	t.Run("error with exactMatch true", func(t *testing.T) {
		result, err := repo.FindByPostalCode(ctx, "01234567", true)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection error")
		assert.ErrorIs(t, err, errMsg.ErrGet)
	})

	t.Run("error with exactMatch false", func(t *testing.T) {
		result, err := repo.FindByPostalCode(ctx, "01234", false)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection error")
		assert.ErrorIs(t, err, errMsg.ErrGet)
	})

	mockDB.AssertExpectations(t)
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
		name  string
		city  string
		state string
	}{
		{"cidade e estado preenchidos", "São Paulo", "SP"},
		{"cidade vazia", "", "SP"},
		{"estado vazio", "São Paulo", ""},
		{"ambos vazios", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := repo.FindByCityAndState(ctx, tc.city, tc.state)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Empty(t, result)
		})
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
	assert.Contains(t, err.Error(), "db error")
	mockDB.AssertExpectations(t)
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
	mockDB.AssertExpectations(t)
}

func TestAddressFilter_FindByPostalCodeV2_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result, err := repo.FindByPostalCodeV2(ctx, "01234567", true)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

// ============================================================================
// Testes do método FindByPostalCodeImproved
// ============================================================================

func TestAddressFilter_FindByPostalCodeImproved_Success(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	t.Run("exact match true", func(t *testing.T) {
		result, err := repo.FindByPostalCodeImproved(ctx, "01234567", true)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("exact match false", func(t *testing.T) {
		result, err := repo.FindByPostalCodeImproved(ctx, "01234-567", false)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}

func TestAddressFilter_FindByPostalCodeImproved_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &addressFilterRepo{db: mockDB}
	ctx := context.Background()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result, err := repo.FindByPostalCodeImproved(ctx, "01234567", true)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

// ============================================================================
// Testes da função auxiliar ensureNonNilSlice
// ============================================================================

func TestEnsureNonNilSlice(t *testing.T) {
	addr1 := &model.Address{ID: 1, Street: "Rua 1"}
	addr2 := &model.Address{ID: 2, Street: "Rua 2"}

	tests := []struct {
		name  string
		input []*model.Address
	}{
		{"slice nil", nil},
		{"slice vazia", []*model.Address{}},
		{"um elemento", []*model.Address{addr1}},
		{"dois elementos", []*model.Address{addr1, addr2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureNonNilAddressSlice(tt.input)

			assert.NotNil(t, result)
			assert.Len(t, result, len(tt.input))

			for i := range tt.input {
				assert.Equal(t, tt.input[i], result[i])
			}
		})
	}
}

// ============================================================================
// Testes de limite/offset
// ============================================================================

func TestAddressFilter_LimitOffset(t *testing.T) {
	cases := []struct {
		name   string
		limit  int
		offset int
	}{
		{"limit 0 offset 0", 0, 0},
		{"limit 10 offset 0", 10, 0},
		{"limit 50 offset 100", 50, 100},
		{"limit 100 offset 1000", 100, 1000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := new(mockDb.MockDatabase)
			repo := &addressFilterRepo{db: mockDB}
			ctx := context.Background()
			isActive := true

			rows := new(mockDb.MockRows)
			rows.On("Next").Return(false)
			rows.On("Err").Return(nil)
			rows.On("Close").Return()

			f := &filterAddress.AddressFilter{
				BaseFilter: filter.BaseFilter{
					Limit:  tc.limit,
					Offset: tc.offset,
				},
				City:     "São Paulo",
				IsActive: &isActive,
			}

			mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

			_, err := repo.Filter(ctx, f)
			assert.NoError(t, err)
		})
	}
}

// ============================================================================
// Testes do mapa de campos de ordenação
// ============================================================================

func TestAddressAllowedSortFields(t *testing.T) {
	// Campos que DEVEM existir no mapa
	expectedFields := []string{
		"id", "user_id", "client_cpf_id", "supplier_id",
		"street", "street_number", "city", "state",
		"country", "postal_code", "is_active", "created_at", "updated_at",
	}

	for _, field := range expectedFields {
		t.Run(field, func(t *testing.T) {
			lowerField := strings.ToLower(field)
			_, exists := addressAllowedSortFields[lowerField]
			assert.True(t, exists, "Campo %s deve estar em addressAllowedSortFields", field)
		})
	}

	// Campos que NÃO devem existir no mapa
	invalidFields := []string{"invalid", "name", "email", "description", "phone", ""}
	for _, field := range invalidFields {
		t.Run("invalid_"+field, func(t *testing.T) {
			_, exists := addressAllowedSortFields[field]
			assert.False(t, exists, "Campo %s não deve estar em addressAllowedSortFields", field)
		})
	}

	t.Run("field normalization", func(t *testing.T) {
		testField := "CrEaTeD_At"
		lowerField := strings.ToLower(testField)
		mappedField, exists := addressAllowedSortFields[lowerField]

		assert.True(t, exists)
		assert.Equal(t, "created_at", mappedField)
	})
}

// ============================================================================
// Helper functions
// ============================================================================

func boolPtr(b bool) *bool {
	return &b
}
