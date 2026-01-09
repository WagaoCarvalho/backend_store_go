package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	filterSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplier_Filter(t *testing.T) {
	t.Run("successfully get all suppliers", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Preparar valores para ponteiros
		cpfValue := "123.456.789-00"
		cnpjValue := "" // string vazia

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // id
			mock.AnythingOfType("*string"),    // name
			mock.AnythingOfType("**string"),   // cpf
			mock.AnythingOfType("**string"),   // cnpj
			mock.AnythingOfType("*string"),    // description
			mock.AnythingOfType("*bool"),      // status
			mock.AnythingOfType("*int"),       // version
			mock.AnythingOfType("*time.Time"), // created_at
			mock.AnythingOfType("*time.Time"), // updated_at
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			*args[1].(*string) = "Fornecedor Teste"
			*args[2].(**string) = &cpfValue
			*args[3].(**string) = &cnpjValue
			*args[4].(*string) = "Descrição do fornecedor"
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now
			*args[8].(*time.Time) = now
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Fornecedor Teste", result[0].Name)
		assert.NotNil(t, result[0].CPF)
		assert.Equal(t, "123.456.789-00", *result[0].CPF)
		assert.NotNil(t, result[0].CNPJ)
		assert.Equal(t, "", *result[0].CNPJ)
		assert.Equal(t, "Descrição do fornecedor", result[0].Description)
		assert.True(t, result[0].Status)
		assert.Equal(t, 1, result[0].Version)
		assert.WithinDuration(t, now, result[0].CreatedAt, time.Second)
		assert.WithinDuration(t, now, result[0].UpdatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter suppliers by name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Preparar valores para ponteiros
		cpfValue := "987.654.321-00"
		cnpjValue := "12.345.678/0001-90"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // id
			mock.AnythingOfType("*string"),    // name
			mock.AnythingOfType("**string"),   // cpf
			mock.AnythingOfType("**string"),   // cnpj
			mock.AnythingOfType("*string"),    // description
			mock.AnythingOfType("*bool"),      // status
			mock.AnythingOfType("*int"),       // version
			mock.AnythingOfType("*time.Time"), // created_at
			mock.AnythingOfType("*time.Time"), // updated_at
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 2
			*args[1].(*string) = "Fornecedor XYZ"
			*args[2].(**string) = &cpfValue
			*args[3].(**string) = &cnpjValue
			*args[4].(*string) = "Outro fornecedor"
			*args[5].(*bool) = true
			*args[6].(*int) = 2
			*args[7].(*time.Time) = now.Add(-24 * time.Hour)
			*args[8].(*time.Time) = now
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "name",
				SortOrder: "desc",
			},
			Name: "XYZ",
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(2), result[0].ID)
		assert.Equal(t, "Fornecedor XYZ", result[0].Name)
		assert.NotNil(t, result[0].CPF)
		assert.Equal(t, "987.654.321-00", *result[0].CPF)
		assert.NotNil(t, result[0].CNPJ)
		assert.Equal(t, "12.345.678/0001-90", *result[0].CNPJ)
		assert.Equal(t, "Outro fornecedor", result[0].Description)
		assert.True(t, result[0].Status)
		assert.Equal(t, 2, result[0].Version)
		assert.WithinDuration(t, now.Add(-24*time.Hour), result[0].CreatedAt, time.Second)
		assert.WithinDuration(t, now, result[0].UpdatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter suppliers by status and date range", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		createdFrom := now.Add(-7 * 24 * time.Hour)
		createdTo := now.Add(-1 * 24 * time.Hour)
		mockRows := new(mockDb.MockRows)

		// Primeiro fornecedor
		cpfValue1 := "111.222.333-44"
		var cnpjValue1 *string // nil para CNPJ

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 3
			*args[1].(*string) = "Fornecedor Inativo"
			*args[2].(**string) = &cpfValue1
			*args[3].(**string) = cnpjValue1 // nil pointer
			*args[4].(*string) = "Fornecedor desativado"
			*args[5].(*bool) = false
			*args[6].(*int) = 3
			*args[7].(*time.Time) = now.Add(-3 * 24 * time.Hour)
			*args[8].(*time.Time) = now.Add(-2 * 24 * time.Hour)
		}).Return(nil).Once()

		// Segundo fornecedor
		var cpfValue2 *string // nil para CPF
		cnpjValue2 := "98.765.432/0001-10"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 4
			*args[1].(*string) = "Fornecedor PJ"
			*args[2].(**string) = cpfValue2 // nil pointer
			*args[3].(**string) = &cnpjValue2
			*args[4].(*string) = "Fornecedor pessoa jurídica"
			*args[5].(*bool) = false
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-5 * 24 * time.Hour)
			*args[8].(*time.Time) = now.Add(-4 * 24 * time.Hour)
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		status := false
		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  20,
				Offset: 0,
			},
			Status:      &status,
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verificar primeiro fornecedor
		assert.Equal(t, int64(3), result[0].ID)
		assert.Equal(t, "Fornecedor Inativo", result[0].Name)
		assert.NotNil(t, result[0].CPF)
		assert.Equal(t, "111.222.333-44", *result[0].CPF)
		assert.Nil(t, result[0].CNPJ) // CNPJ deve ser nil
		assert.Equal(t, "Fornecedor desativado", result[0].Description)
		assert.False(t, result[0].Status)
		assert.Equal(t, 3, result[0].Version)

		// Verificar segundo fornecedor
		assert.Equal(t, int64(4), result[1].ID)
		assert.Equal(t, "Fornecedor PJ", result[1].Name)
		assert.Nil(t, result[1].CPF) // CPF deve ser nil
		assert.NotNil(t, result[1].CNPJ)
		assert.Equal(t, "98.765.432/0001-10", *result[1].CNPJ)
		assert.Equal(t, "Fornecedor pessoa jurídica", result[1].Description)
		assert.False(t, result[1].Status)
		assert.Equal(t, 1, result[1].Version)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter suppliers by CPF and CNPJ", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Mock para um único fornecedor com CPF específico
		cpfValue := "123.456.789-00"
		cnpjValue := "12.345.678/0001-90"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			*args[1].(*string) = "Fornecedor Específico"
			*args[2].(**string) = &cpfValue
			*args[3].(**string) = &cnpjValue
			*args[4].(*string) = "Fornecedor com CPF e CNPJ específicos"
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-24 * time.Hour)
			*args[8].(*time.Time) = now
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// Criar filtro com CPF e CNPJ específicos
		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			CPF:  "123.456.789-00",
			CNPJ: "12.345.678/0001-90",
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Fornecedor Específico", result[0].Name)
		assert.NotNil(t, result[0].CPF)
		assert.Equal(t, "123.456.789-00", *result[0].CPF)
		assert.NotNil(t, result[0].CNPJ)
		assert.Equal(t, "12.345.678/0001-90", *result[0].CNPJ)
		assert.Equal(t, "Fornecedor com CPF e CNPJ específicos", result[0].Description)
		assert.True(t, result[0].Status)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter suppliers by updated_at range", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		updatedFrom := now.Add(-72 * time.Hour) // 3 dias atrás
		updatedTo := now.Add(-24 * time.Hour)   // 1 dia atrás

		mockRows := new(mockDb.MockRows)

		// Mock para dois fornecedores dentro do intervalo de updated_at
		// Primeiro fornecedor
		cpfValue1 := "111.222.333-44"
		cnpjValue1 := "11.222.333/0001-44"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			*args[1].(*string) = "Fornecedor Atualizado Recente"
			*args[2].(**string) = &cpfValue1
			*args[3].(**string) = &cnpjValue1
			*args[4].(*string) = "Atualizado há 2 dias"
			*args[5].(*bool) = true
			*args[6].(*int) = 3
			*args[7].(*time.Time) = now.Add(-10 * 24 * time.Hour) // created_at antigo
			*args[8].(*time.Time) = now.Add(-48 * time.Hour)      // updated_at dentro do range
		}).Return(nil).Once()

		// Segundo fornecedor
		cpfValue2 := "555.666.777-88"
		var cnpjValue2 *string // nil para CNPJ

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 2
			*args[1].(*string) = "Fornecedor Atualizado Limite"
			*args[2].(**string) = &cpfValue2
			*args[3].(**string) = cnpjValue2 // nil
			*args[4].(*string) = "Atualizado há exatamente 3 dias"
			*args[5].(*bool) = false
			*args[6].(*int) = 2
			*args[7].(*time.Time) = now.Add(-15 * 24 * time.Hour) // created_at muito antigo
			*args[8].(*time.Time) = now.Add(-72 * time.Hour)      // updated_at no limite inferior
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// Criar filtro com range de updated_at
		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "updated_at",
				SortOrder: "desc",
			},
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verificar ordenação por updated_at desc
		assert.Equal(t, int64(1), result[0].ID) // updated_at mais recente primeiro
		assert.Equal(t, "Fornecedor Atualizado Recente", result[0].Name)
		assert.WithinDuration(t, now.Add(-48*time.Hour), result[0].UpdatedAt, time.Second)

		assert.Equal(t, int64(2), result[1].ID)
		assert.Equal(t, "Fornecedor Atualizado Limite", result[1].Name)
		assert.Nil(t, result[1].CNPJ) // CNPJ deve ser nil
		assert.WithinDuration(t, now.Add(-72*time.Hour), result[1].UpdatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter suppliers with invalid sort order defaults to asc", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Mock para dois fornecedores
		// Primeiro fornecedor (mais antigo)
		cpfValue1 := "111.111.111-11"
		cnpjValue1 := "11.111.111/0001-11"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			*args[1].(*string) = "Fornecedor A"
			*args[2].(**string) = &cpfValue1
			*args[3].(**string) = &cnpjValue1
			*args[4].(*string) = "Mais antigo"
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-48 * time.Hour) // Mais antigo
			*args[8].(*time.Time) = now.Add(-48 * time.Hour)
		}).Return(nil).Once()

		// Segundo fornecedor (mais recente)
		cpfValue2 := "222.222.222-22"
		cnpjValue2 := "22.222.222/0001-22"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 2
			*args[1].(*string) = "Fornecedor B"
			*args[2].(**string) = &cpfValue2
			*args[3].(**string) = &cnpjValue2
			*args[4].(*string) = "Mais recente"
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-24 * time.Hour) // Mais recente
			*args[8].(*time.Time) = now.Add(-24 * time.Hour)
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// Criar filtro com sortOrder inválido (deve default para "asc")
		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "created_at",
				SortOrder: "INVALID_ORDER", // Ordem inválida
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verificar que ordenação é ASC (mais antigo primeiro)
		assert.Equal(t, int64(1), result[0].ID) // Fornecedor A (mais antigo)
		assert.Equal(t, "Fornecedor A", result[0].Name)
		assert.WithinDuration(t, now.Add(-48*time.Hour), result[0].CreatedAt, time.Second)

		assert.Equal(t, int64(2), result[1].ID) // Fornecedor B (mais recente)
		assert.Equal(t, "Fornecedor B", result[1].Name)
		assert.WithinDuration(t, now.Add(-24*time.Hour), result[1].CreatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("returns error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		// Mock para retornar erro na query
		expectedErr := fmt.Errorf("database connection failed")

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(nil, expectedErr)

		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		result, err := repo.Filter(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection failed")
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("returns error when scanning row fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(fmt.Errorf("scan error: invalid data type")) // Retorna erro no Scan

		mockRows.On("Close").Return()

		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "scan error: invalid data type")
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("returns error when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)

		// Primeira linha OK
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			*args[1].(*string) = "Fornecedor Teste"
			cpfValue := "111.222.333-44"
			*args[2].(**string) = &cpfValue
			var cnpjValue *string
			*args[3].(**string) = cnpjValue
			*args[4].(*string) = "Descrição"
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = time.Now()
			*args[8].(*time.Time) = time.Now()
		}).Return(nil).Once()

		// Erro na próxima chamada de Next
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(fmt.Errorf("cursor error: connection lost"))
		mockRows.On("Close").Return()

		filter := &filterSupplier.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cursor error: connection lost")
		assert.Contains(t, err.Error(), errMsg.ErrIterate.Error())

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
