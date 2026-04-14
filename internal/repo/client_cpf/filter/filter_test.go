package repo

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	filterClient "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// Testes do método Filter
// ============================================================================

func TestClientFilter_Filter_Success_WithResults(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()
	now := time.Now()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		if id, ok := args[0].(*int64); ok {
			*id = 1
		}
		if name, ok := args[1].(*string); ok {
			*name = "João Silva"
		}
		if email, ok := args[2].(*string); ok {
			*email = "joao@email.com"
		}
		if cpf, ok := args[3].(*string); ok {
			*cpf = "12345678909"
		}
		if desc, ok := args[4].(*string); ok {
			*desc = "Cliente teste"
		}
		if status, ok := args[5].(*bool); ok {
			*status = true
		}
		if version, ok := args[6].(*int); ok {
			*version = 1
		}
		if createdAt, ok := args[7].(*time.Time); ok {
			*createdAt = now
		}
		if updatedAt, ok := args[8].(*time.Time); ok {
			*updatedAt = now
		}
	}).Return(nil).Once()
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "João Silva", result[0].Name)
	assert.Equal(t, "joao@email.com", result[0].Email)
	assert.Equal(t, "12345678909", result[0].CPF)
	assert.True(t, result[0].Status)
	assert.Equal(t, 1, result[0].Version)
	mockDB.AssertExpectations(t)
	rows.AssertExpectations(t)
}

func TestClientFilter_Filter_WithAllFilters(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	status := true
	version := 1
	createdFrom := time.Now().Add(-24 * time.Hour)
	createdTo := time.Now()
	updatedFrom := time.Now().Add(-12 * time.Hour)
	updatedTo := time.Now()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "name",
			SortOrder: "asc",
			Limit:     10,
			Offset:    0,
		},
		Name:        "Teste",
		Email:       "teste@email.com",
		CPF:         "11122233344",
		Status:      &status,
		Version:     &version,
		CreatedFrom: &createdFrom,
		CreatedTo:   &createdTo,
		UpdatedFrom: &updatedFrom,
		UpdatedTo:   &updatedTo,
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockDB.AssertExpectations(t)
}

func TestClientFilter_Filter_EmptyResult(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockDB.AssertExpectations(t)
}

func TestClientFilter_Filter_MultipleRows(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)

	// Primeira linha
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		if id, ok := args[0].(*int64); ok {
			*id = 1
		}
		if name, ok := args[1].(*string); ok {
			*name = "Cliente 1"
		}
	}).Return(nil).Once()

	// Segunda linha
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		if id, ok := args[0].(*int64); ok {
			*id = 2
		}
		if name, ok := args[1].(*string); ok {
			*name = "Cliente 2"
		}
	}).Return(nil).Once()

	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
	}

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Cliente 1", result[0].Name)
	assert.Equal(t, "Cliente 2", result[1].Name)
	mockDB.AssertExpectations(t)
}

// ============================================================================
// Testes de ordenação
// ============================================================================

func TestClientFilter_SortField_Valid(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "name",
			SortOrder: "asc",
			Limit:     10,
		},
	}

	mockDB.
		On("Query",
			ctx,
			mock.MatchedBy(func(q string) bool {
				return strings.Contains(q, "ORDER BY name asc")
			}),
			mock.Anything,
		).
		Return(rows, nil)

	_, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
}

func TestClientFilter_SortField_Invalid(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "invalid_field",
			SortOrder: "asc",
			Limit:     10,
		},
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

func TestClientFilter_SortOrder_Invalid(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "created_at",
			SortOrder: "INVALID",
			Limit:     10,
		},
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

func TestClientFilter_SortOrder_Descending(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{
			SortBy:    "name",
			SortOrder: "desc",
			Limit:     10,
		},
	}

	mockDB.
		On("Query",
			ctx,
			mock.MatchedBy(func(q string) bool {
				return strings.Contains(q, "ORDER BY name desc")
			}),
			mock.Anything,
		).
		Return(rows, nil)

	_, err := repo.Filter(ctx, f)
	assert.NoError(t, err)
}

func TestClientFilter_AllSortFields(t *testing.T) {
	for sortField := range clientCpfAllowedSortFields {
		t.Run(sortField, func(t *testing.T) {
			mockDB := new(mockDb.MockDatabase)
			repo := &clientCpfFilterRepo{db: mockDB}
			ctx := context.Background()

			rows := new(mockDb.MockRows)
			rows.On("Next").Return(false)
			rows.On("Err").Return(nil)
			rows.On("Close").Return()

			f := &filterClient.ClientCpfFilter{
				BaseFilter: filter.BaseFilter{
					SortBy:    sortField,
					SortOrder: "asc",
					Limit:     10,
				},
			}

			mockDB.
				On("Query",
					ctx,
					mock.MatchedBy(func(q string) bool {
						expectedField := clientCpfAllowedSortFields[sortField]
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

func TestClientFilter_CaseInsensitive(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	t.Run("sort field uppercase", func(t *testing.T) {
		f := &filterClient.ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				SortBy:    "CREATED_AT",
				SortOrder: "asc",
				Limit:     10,
			},
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
		f := &filterClient.ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				SortBy:    "name",
				SortOrder: "DESC",
				Limit:     10,
			},
		}

		mockDB.
			On("Query",
				ctx,
				mock.MatchedBy(func(q string) bool {
					return strings.Contains(q, "ORDER BY name desc")
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

func TestClientFilter_QueryError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	dbErr := errors.New("db error")
	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
	}

	result, err := repo.Filter(ctx, f)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errMsg.ErrGet)
	assert.Contains(t, err.Error(), "db error")
}

func TestClientFilter_ScanError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(true).Once()
	rows.On("Scan",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything,
	).Return(errors.New("scan error"))
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
	}

	result, err := repo.Filter(ctx, f)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errMsg.ErrScan)
	assert.Contains(t, err.Error(), "scan error")
}

func TestClientFilter_RowsError(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(errors.New("iterate error"))
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	f := &filterClient.ClientCpfFilter{
		BaseFilter: filter.BaseFilter{Limit: 10},
	}

	result, err := repo.Filter(ctx, f)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errMsg.ErrIterate)
	assert.Contains(t, err.Error(), "iterate error")
}

// ============================================================================
// Testes de métodos auxiliares
// ============================================================================

func TestClientFilter_FindActive(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	filter := &filterClient.ClientCpfFilter{
		Name: "Teste",
	}

	result, err := repo.FindActive(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.NotNil(t, filter.Status)
	assert.True(t, *filter.Status)
}

// ============================================================================
// Testes do método FindByCPF (versão 1)
// ============================================================================

func TestClientFilter_FindByCPF_Success(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	t.Run("exact match true", func(t *testing.T) {
		result, err := repo.FindByCPF(ctx, "12345678909", true)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("exact match false", func(t *testing.T) {
		result, err := repo.FindByCPF(ctx, "123", false)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}

// Teste específico para cobrir o trecho de erro no FindByCPF
func TestClientFilter_FindByCPF_Error(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	// Configura o mock para retornar erro
	dbErr := errors.New("database connection error")
	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

	// Testa com exactMatch true
	t.Run("error with exactMatch true", func(t *testing.T) {
		result, err := repo.FindByCPF(ctx, "12345678909", true)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection error")
		assert.ErrorIs(t, err, errMsg.ErrGet)
	})

	// Testa com exactMatch false
	t.Run("error with exactMatch false", func(t *testing.T) {
		result, err := repo.FindByCPF(ctx, "123", false)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection error")
		assert.ErrorIs(t, err, errMsg.ErrGet)
	})

	mockDB.AssertExpectations(t)
}

// Versão mais compacta do teste de erro
func TestClientFilter_FindByCPF_ErrorCoverage(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	// Mock de erro
	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("forced error for coverage"))

	// Testa apenas um caso (já cobre o trecho if err != nil)
	result, err := repo.FindByCPF(ctx, "12345678909", true)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestClientFilter_FindByName(t *testing.T) {
	mockDB := new(mockDb.MockDatabase)
	repo := &clientCpfFilterRepo{db: mockDB}
	ctx := context.Background()

	rows := new(mockDb.MockRows)
	rows.On("Next").Return(false)
	rows.On("Err").Return(nil)
	rows.On("Close").Return()

	mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

	result, err := repo.FindByName(ctx, "João")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// ============================================================================
// Testes da função auxiliar ensureNonNilClientSlice
// ============================================================================

func TestEnsureNonNilClientSlice(t *testing.T) {
	client1 := &model.ClientCpf{ID: 1, Name: "Cliente 1"}
	client2 := &model.ClientCpf{ID: 2, Name: "Cliente 2"}

	tests := []struct {
		name  string
		input []*model.ClientCpf
	}{
		{"slice nil", nil},
		{"slice vazia", []*model.ClientCpf{}},
		{"um elemento", []*model.ClientCpf{client1}},
		{"dois elementos", []*model.ClientCpf{client1, client2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureNonNilClientSlice(tt.input)

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

func TestClientFilter_LimitOffset(t *testing.T) {
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
			repo := &clientCpfFilterRepo{db: mockDB}
			ctx := context.Background()

			rows := new(mockDb.MockRows)
			rows.On("Next").Return(false)
			rows.On("Err").Return(nil)
			rows.On("Close").Return()

			f := &filterClient.ClientCpfFilter{
				BaseFilter: filter.BaseFilter{
					Limit:  tc.limit,
					Offset: tc.offset,
				},
				Name: "Teste",
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

// ============================================================================
// Testes do mapa de campos de ordenação
// ============================================================================

func TestClientAllowedSortFields(t *testing.T) {
	// Campos que DEVEM existir no mapa
	expectedFields := []string{
		"id", "name", "email", "cpf", "description", "status", "version", "created_at", "updated_at",
	}

	for _, field := range expectedFields {
		t.Run(field, func(t *testing.T) {
			lowerField := strings.ToLower(field)
			_, exists := clientCpfAllowedSortFields[lowerField]
			assert.True(t, exists, "Campo %s deve estar em clientCpfAllowedSortFields", field)
		})
	}

	// Campos que NÃO devem existir no mapa
	invalidFields := []string{"invalid", "phone", "address", "city", ""}
	for _, field := range invalidFields {
		t.Run("invalid_"+field, func(t *testing.T) {
			_, exists := clientCpfAllowedSortFields[field]
			assert.False(t, exists, "Campo %s não deve estar em clientCpfAllowedSortFields", field)
		})
	}

	t.Run("field normalization", func(t *testing.T) {
		testField := "CrEaTeD_At"
		lowerField := strings.ToLower(testField)
		mappedField, exists := clientCpfAllowedSortFields[lowerField]

		assert.True(t, exists)
		assert.Equal(t, "created_at", mappedField)
	})
}
