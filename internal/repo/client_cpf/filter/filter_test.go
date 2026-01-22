package repo

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	filterClient "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_Filter(t *testing.T) {

	t.Run("successfully get all clients", func(t *testing.T) {
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
				*cpf = "123.456.789-09"
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

		rows.On("Next").Return(false).Once()
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
		mockDB.AssertExpectations(t)
		rows.AssertExpectations(t)
	})

	t.Run("with all filters applied", func(t *testing.T) {
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
			CPF:         "111.222.333-44",
			Status:      &status,
			Version:     &version,
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
		}

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

		_, err := repo.Filter(ctx, f)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("empty result", func(t *testing.T) {
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
		assert.Empty(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("multiple rows returned", func(t *testing.T) {
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
	})

	t.Run("uses allowed sort field", func(t *testing.T) {
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
	})

	t.Run("defaults sort field to created_at when invalid", func(t *testing.T) {
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
	})

	t.Run("defaults sort order to asc when invalid", func(t *testing.T) {
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
	})

	t.Run("sort order descending", func(t *testing.T) {
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
	})

	t.Run("test all allowed sort fields", func(t *testing.T) {

		for sortField := range allowedSortFields {
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
							expectedField := allowedSortFields[sortField]
							return strings.Contains(q, "ORDER BY "+expectedField+" asc")
						}),
						mock.Anything,
					).
					Return(rows, nil)

				_, err := repo.Filter(ctx, f)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
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
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
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
	})

	t.Run("return ErrIterate when rows.Err fails", func(t *testing.T) {
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
	})

	t.Run("test case insensitive sort field", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfFilterRepo{db: mockDB}
		ctx := context.Background()

		rows := new(mockDb.MockRows)
		rows.On("Next").Return(false)
		rows.On("Err").Return(nil)
		rows.On("Close").Return()

		// Testando com letras maiúsculas
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

	t.Run("test case insensitive sort order", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfFilterRepo{db: mockDB}
		ctx := context.Background()

		rows := new(mockDb.MockRows)
		rows.On("Next").Return(false)
		rows.On("Err").Return(nil)
		rows.On("Close").Return()

		// Testando com letras maiúsculas
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

	t.Run("test with only some filters", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfFilterRepo{db: mockDB}
		ctx := context.Background()

		rows := new(mockDb.MockRows)
		rows.On("Next").Return(false)
		rows.On("Err").Return(nil)
		rows.On("Close").Return()

		// Testando apenas com alguns filtros
		f := &filterClient.ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				Limit: 10,
			},
			Name:  "Partial",
			Email: "partial@email.com",
			// CPF, Status, etc. não preenchidos
		}

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

		_, err := repo.Filter(ctx, f)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}

func TestClient_Filter_CompleteCoverage(t *testing.T) {

	cases := []struct {
		name  string
		limit int
	}{
		{"limit 0", 0},
		{"limit 1", 1},
		{"limit 50", 50},
		{"limit 100", 100},
		{"limit 101", 101},
		{"limit 200", 200},
		{"limit negative", -1},
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
					Limit: tc.limit,
				},
				Name: "Teste",
			}

			mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(rows, nil)

			_, err := repo.Filter(ctx, f)
			assert.NoError(t, err)
		})
	}
}

func TestAllowedSortFields(t *testing.T) {
	// Testa se o mapa allowedSortFields contém todos os campos esperados
	expectedFields := []string{
		"id", "name", "email", "cpf", "status", "version", "created_at", "updated_at",
	}

	for _, field := range expectedFields {
		t.Run(field, func(t *testing.T) {
			lowerField := strings.ToLower(field)
			_, exists := allowedSortFields[lowerField]
			assert.True(t, exists, "Campo %s deve estar em allowedSortFields", field)
		})
	}

	// Testa que campos não mapeados retornam o próprio valor em minúsculas
	t.Run("field normalization", func(t *testing.T) {
		testField := "CrEaTeD_At"
		lowerField := strings.ToLower(testField)
		mappedField, exists := allowedSortFields[lowerField]

		assert.True(t, exists)
		assert.Equal(t, "created_at", mappedField)
	})
}
