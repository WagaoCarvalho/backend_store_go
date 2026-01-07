package repo

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	baseFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductFilterRepo_Filter(t *testing.T) {
	t.Run("successfully get all products", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Usando mock.Anything para ser mais flexível
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				// Assumindo a ordem baseada no seu código original
				if ptr, ok := args[0].(*int64); ok {
					*ptr = 1
				}
				if ptr, ok := args[1].(**int64); ok {
					supplierID := int64(10)
					*ptr = &supplierID
				}
				if ptr, ok := args[2].(*string); ok {
					*ptr = "Produto Teste"
				}
				if ptr, ok := args[3].(**string); ok {
					manufacturer := "Fabricante Teste"
					*ptr = &manufacturer
				}
				if ptr, ok := args[4].(**string); ok {
					desc := "Descrição do produto"
					*ptr = &desc
				}
				if ptr, ok := args[5].(*float64); ok {
					*ptr = 50.0
				}
				if ptr, ok := args[6].(*float64); ok {
					*ptr = 100.0
				}
				if ptr, ok := args[7].(*int); ok {
					*ptr = 100
				}
				if ptr, ok := args[8].(*int); ok {
					*ptr = 10
				}
				if ptr, ok := args[9].(**int); ok {
					maxStock := 500
					*ptr = &maxStock
				}
				if ptr, ok := args[10].(**string); ok {
					barcode := "1234567890123"
					*ptr = &barcode
				}
				if ptr, ok := args[11].(*bool); ok {
					*ptr = true
				}
				if ptr, ok := args[12].(*int); ok {
					*ptr = 1
				}
				if ptr, ok := args[13].(*bool); ok {
					*ptr = true
				}
				if ptr, ok := args[14].(*float64); ok {
					*ptr = 0.0
				}
				if ptr, ok := args[15].(*float64); ok {
					*ptr = 30.0
				}
				if ptr, ok := args[16].(*time.Time); ok {
					*ptr = now
				}
				if ptr, ok := args[17].(*time.Time); ok {
					*ptr = now
				}
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters price ranges correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Configura o mock para retornar um produto
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				// Simula um produto com preços específicos
				if ptr, ok := args[0].(*int64); ok {
					*ptr = 2
				}
				if ptr, ok := args[1].(**int64); ok {
					supplierID := int64(15)
					*ptr = &supplierID
				}
				if ptr, ok := args[2].(*string); ok {
					*ptr = "Produto Filtrado"
				}
				if ptr, ok := args[3].(**string); ok {
					manufacturer := "Marca Premium"
					*ptr = &manufacturer
				}
				if ptr, ok := args[4].(**string); ok {
					desc := "Produto com preço entre 100 e 200"
					*ptr = &desc
				}
				if ptr, ok := args[5].(*float64); ok { // cost_price
					*ptr = 80.0
				}
				if ptr, ok := args[6].(*float64); ok { // sale_price
					*ptr = 150.0
				}
				if ptr, ok := args[7].(*int); ok { // stock_quantity
					*ptr = 75
				}
				if ptr, ok := args[8].(*int); ok { // min_stock
					*ptr = 5
				}
				if ptr, ok := args[9].(**int); ok { // max_stock
					maxStock := 200
					*ptr = &maxStock
				}
				if ptr, ok := args[10].(**string); ok { // barcode
					barcode := "5555555555555"
					*ptr = &barcode
				}
				if ptr, ok := args[11].(*bool); ok { // status
					*ptr = true
				}
				if ptr, ok := args[12].(*int); ok { // version
					*ptr = 1
				}
				if ptr, ok := args[13].(*bool); ok { // allow_discount
					*ptr = true
				}
				if ptr, ok := args[14].(*float64); ok { // min_discount_percent
					*ptr = 5.0
				}
				if ptr, ok := args[15].(*float64); ok { // max_discount_percent
					*ptr = 20.0
				}
				if ptr, ok := args[16].(*time.Time); ok { // created_at
					*ptr = now
				}
				if ptr, ok := args[17].(*time.Time); ok { // updated_at
					*ptr = now
				}
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// Define os filtros de preço
		minCostPrice := 50.0
		maxCostPrice := 100.0
		minSalePrice := 100.0
		maxSalePrice := 200.0

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			MinCostPrice: &minCostPrice,
			MaxCostPrice: &maxCostPrice,
			MinSalePrice: &minSalePrice,
			MaxSalePrice: &maxSalePrice,
		}

		// Mock da query verificando se os filtros são aplicados corretamente
		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			// Verifica se os argumentos contêm os valores dos filtros
			// A ordem depende de como sua função constrói a query

			// Conta quantos argumentos existem
			if len(args) < 4 {
				return false
			}

			// Encontra os valores dos filtros nos argumentos
			// Nota: a ordem exata depende da construção da query no seu código
			// Vamos assumir que os filtros são adicionados na ordem:
			// 1. MinCostPrice, 2. MaxCostPrice, 3. MinSalePrice, 4. MaxSalePrice

			foundMinCost := false
			foundMaxCost := false
			foundMinSale := false
			foundMaxSale := false

			for _, arg := range args {
				switch v := arg.(type) {
				case float64:
					if v == 50.0 {
						foundMinCost = true
					} else if v == 100.0 {
						// Pode ser tanto maxCostPrice quanto minSalePrice
						// Vamos verificar pelo contexto
						foundMaxCost = true
						foundMinSale = true
					} else if v == 200.0 {
						foundMaxSale = true
					}
				}
			}

			return foundMinCost && foundMaxCost && foundMinSale && foundMaxSale
		})).Return(mockRows, nil)

		// Executa o método
		result, err := repo.Filter(ctx, filter)

		// Verificações
		assert.NoError(t, err)
		assert.Len(t, result, 1)

		// Verifica se o produto retornado está dentro dos filtros
		product := result[0]
		assert.Equal(t, int64(2), product.ID)
		assert.Equal(t, "Produto Filtrado", product.ProductName)
		assert.Equal(t, 80.0, product.CostPrice)
		assert.Equal(t, 150.0, product.SalePrice)

		// Verifica se os preços estão dentro dos intervalos filtrados
		assert.GreaterOrEqual(t, product.CostPrice, 50.0)
		assert.LessOrEqual(t, product.CostPrice, 100.0)
		assert.GreaterOrEqual(t, product.SalePrice, 100.0)
		assert.LessOrEqual(t, product.SalePrice, 200.0)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters product_name, manufacturer and stock correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Configura o mock para retornar um produto
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				// Simula um produto que atende aos filtros
				if ptr, ok := args[0].(*int64); ok {
					*ptr = 3
				}
				if ptr, ok := args[1].(**int64); ok {
					supplierID := int64(25)
					*ptr = &supplierID
				}
				if ptr, ok := args[2].(*string); ok {
					*ptr = "Notebook Dell XPS"
				}
				if ptr, ok := args[3].(*string); ok { // Manufacturer é string, não *string
					*ptr = "Dell"
				}
				if ptr, ok := args[4].(*string); ok { // Description é string, não *string
					*ptr = "Notebook de alta performance"
				}
				if ptr, ok := args[5].(*float64); ok { // cost_price
					*ptr = 4000.0
				}
				if ptr, ok := args[6].(*float64); ok { // sale_price
					*ptr = 6000.0
				}
				if ptr, ok := args[7].(*int); ok { // stock_quantity
					*ptr = 45
				}
				if ptr, ok := args[8].(*int); ok { // min_stock
					*ptr = 5
				}
				if ptr, ok := args[9].(**int); ok { // max_stock
					maxStock := 100
					*ptr = &maxStock
				}
				if ptr, ok := args[10].(*string); ok { // barcode é string, não *string
					*ptr = "DELLXPS123456"
				}
				if ptr, ok := args[11].(*bool); ok { // status
					*ptr = true
				}
				if ptr, ok := args[12].(*int); ok { // version
					*ptr = 2
				}
				if ptr, ok := args[13].(*bool); ok { // allow_discount
					*ptr = false
				}
				if ptr, ok := args[14].(*float64); ok { // min_discount_percent
					*ptr = 0.0
				}
				if ptr, ok := args[15].(*float64); ok { // max_discount_percent
					*ptr = 10.0
				}
				if ptr, ok := args[16].(*time.Time); ok { // created_at
					*ptr = now
				}
				if ptr, ok := args[17].(*time.Time); ok { // updated_at
					*ptr = now
				}
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		// Define os filtros: nome contém "Dell", fabricante "Dell", estoque entre 40 e 50
		minStock := 40
		maxStock := 50

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			ProductName:      "Dell",
			Manufacturer:     "Dell",
			MinStockQuantity: &minStock,
			MaxStockQuantity: &maxStock,
		}

		// Mock da query verificando se os filtros são aplicados corretamente
		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			// Verifica se os argumentos contêm os valores dos filtros
			if len(args) < 4 {
				return false
			}

			// Verifica se os argumentos correspondem aos filtros
			hasProductName := false
			hasManufacturer := false
			hasMinStock := false
			hasMaxStock := false

			for _, arg := range args {
				switch v := arg.(type) {
				case string:
					if v == "Dell" {
						// Pode ser product_name ou manufacturer
						// Como ambos são "Dell" neste teste, vamos considerar que ambos estão presentes
						hasProductName = true
						hasManufacturer = true
					}
				case int:
					if v == 40 {
						hasMinStock = true
					} else if v == 50 {
						hasMaxStock = true
					}
				}
			}

			// Retorna true se todos os filtros foram encontrados
			return hasProductName && hasManufacturer && hasMinStock && hasMaxStock
		})).Return(mockRows, nil)

		// Executa o método
		result, err := repo.Filter(ctx, filter)

		// Verificações
		assert.NoError(t, err)
		assert.Len(t, result, 1)

		// Verifica se o produto retornado atende aos filtros
		product := result[0]
		assert.Equal(t, int64(3), product.ID)
		assert.Equal(t, "Notebook Dell XPS", product.ProductName)

		// Verifica se o nome contém "Dell"
		assert.Contains(t, product.ProductName, "Dell")

		// Verifica se o fabricante é "Dell" (agora é string direto, não ponteiro)
		assert.Equal(t, "Dell", product.Manufacturer)

		// Verifica se o estoque está dentro do intervalo
		assert.GreaterOrEqual(t, product.StockQuantity, 40)
		assert.LessOrEqual(t, product.StockQuantity, 50)
		assert.Equal(t, 45, product.StockQuantity)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		// Simula um erro de banco de dados
		dbErr := errors.New("database connection failed")

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		// Mock retorna erro na query
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(nil, dbErr)

		// Executa o método
		result, err := repo.Filter(ctx, filter)

		// Verificações
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		// Simula um erro no scan
		scanErr := errors.New("failed to scan row")

		mockRows := new(mockDb.MockRows)

		// Mock retorna true para Next() (há uma linha)
		mockRows.On("Next").Return(true).Once()

		// Mock do Scan retorna erro
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(scanErr).Once()

		// Garante que Close é chamado
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		// Mock da query retorna rows
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		// Executa o método
		result, err := repo.Filter(ctx, filter)

		// Verificações
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.ErrorContains(t, err, scanErr.Error())

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		// Simula um erro na iteração das linhas
		rowsErr := errors.New("iteration error")

		mockRows := new(mockDb.MockRows)

		// Mock retorna false para Next() (não há linhas)
		mockRows.On("Next").Return(false).Once()

		// Mock do Err retorna erro
		mockRows.On("Err").Return(rowsErr)

		// Garante que Close é chamado
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		// Mock da query retorna rows
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		// Executa o método
		result, err := repo.Filter(ctx, filter)

		// Verificações
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.ErrorContains(t, err, rowsErr.Error())

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("uses allowed sort field when SortBy is valid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				// Preenche valores básicos
				if ptr, ok := args[0].(*int64); ok {
					*ptr = 1
				}
				if ptr, ok := args[2].(*string); ok {
					*ptr = "Produto Teste"
				}
				if ptr, ok := args[6].(*float64); ok {
					*ptr = 100.0
				}
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "sale_price", // Campo válido
				SortOrder: "desc",
			},
		}

		// Verifica se o ORDER BY está correto
		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(strings.ToLower(query), "order by sale_price desc")
		}), mock.Anything).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("defaults to created_at when SortBy is invalid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "invalid_field", // Campo inválido
				SortOrder: "asc",
			},
		}

		// Verifica se usa created_at como padrão
		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(strings.ToLower(query), "order by created_at asc")
		}), mock.Anything).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("defaults sort order to asc when SortOrder is invalid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "product_name",
				SortOrder: "INVALID_ORDER", // Ordenação inválida
			},
		}

		// Verifica se usa asc como padrão
		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(strings.ToLower(query), "order by product_name asc")
		}), mock.Anything).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters supplier_id, status and allow_discount correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				if ptr, ok := args[0].(*int64); ok {
					*ptr = 5
				}
				if ptr, ok := args[2].(*string); ok {
					*ptr = "Produto com Fornecedor"
				}
				if ptr, ok := args[11].(*bool); ok {
					*ptr = true // status
				}
				if ptr, ok := args[13].(*bool); ok {
					*ptr = false // allow_discount
				}
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		supplierID := int64(100)
		status := true
		allowDiscount := false

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			SupplierID:    &supplierID,
			Status:        &status,
			AllowDiscount: &allowDiscount,
		}

		// Verifica se os filtros estão presentes nos argumentos
		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			hasSupplierID := false
			hasStatus := false
			hasAllowDiscount := false

			for _, arg := range args {
				switch v := arg.(type) {
				case int64:
					if v == 100 {
						hasSupplierID = true
					}
				case bool:
					if v == true {
						hasStatus = true
					} else if v == false {
						hasAllowDiscount = true
					}
				}
			}

			return hasSupplierID && hasStatus && hasAllowDiscount
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters barcode and version correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				if ptr, ok := args[0].(*int64); ok {
					*ptr = 6
				}
				if ptr, ok := args[2].(*string); ok {
					*ptr = "Produto por Código"
				}
				if ptr, ok := args[10].(*string); ok {
					*ptr = "7891234567890"
				}
				if ptr, ok := args[12].(*int); ok {
					*ptr = 3 // version
				}
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		barcode := "7891234567890"
		version := 3

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			Barcode: barcode,
			Version: &version,
		}

		// Verifica se os filtros estão presentes
		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			hasBarcode := false
			hasVersion := false

			for _, arg := range args {
				switch v := arg.(type) {
				case string:
					if v == "7891234567890" {
						hasBarcode = true
					}
				case int:
					if v == 3 {
						hasVersion = true
					}
				}
			}

			return hasBarcode && hasVersion
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters date ranges correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				if ptr, ok := args[0].(*int64); ok {
					*ptr = 7
				}
				if ptr, ok := args[2].(*string); ok {
					*ptr = "Produto por Data"
				}
				if ptr, ok := args[16].(*time.Time); ok {
					*ptr = now.Add(-24 * time.Hour) // created_at
				}
				if ptr, ok := args[17].(*time.Time); ok {
					*ptr = now // updated_at
				}
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		createdFrom := now.Add(-48 * time.Hour)
		createdTo := now.Add(-12 * time.Hour)
		updatedFrom := now.Add(-6 * time.Hour)
		updatedTo := now.Add(1 * time.Hour)

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
		}

		// Verifica se os filtros de data estão presentes
		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			hasCreatedFrom := false
			hasCreatedTo := false
			hasUpdatedFrom := false
			hasUpdatedTo := false

			for _, arg := range args {
				switch v := arg.(type) {
				case time.Time:
					if v.Equal(createdFrom) {
						hasCreatedFrom = true
					} else if v.Equal(createdTo) {
						hasCreatedTo = true
					} else if v.Equal(updatedFrom) {
						hasUpdatedFrom = true
					} else if v.Equal(updatedTo) {
						hasUpdatedTo = true
					}
				}
			}

			return hasCreatedFrom && hasCreatedTo && hasUpdatedFrom && hasUpdatedTo
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
