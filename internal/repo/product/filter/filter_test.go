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
	t.Run("error when filter is nil", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		result, err := repo.Filter(ctx, nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockDB.AssertNotCalled(t, "Query")
	})

	t.Run("successfully get all products", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.Anything, // ID
			mock.Anything, // SupplierID
			mock.Anything, // ProductName
			mock.Anything, // Manufacturer
			mock.Anything, // Description
			mock.Anything, // CostPrice
			mock.Anything, // SalePrice
			mock.Anything, // StockQuantity
			mock.Anything, // MinStock
			mock.Anything, // MaxStock
			mock.Anything, // Barcode
			mock.Anything, // Status
			mock.Anything, // Version
			mock.Anything, // AllowDiscount
			mock.Anything, // MinDiscountPercent
			mock.Anything, // MaxDiscountPercent
			mock.Anything, // CreatedAt
			mock.Anything, // UpdatedAt
		).Run(func(args mock.Arguments) {
			// ID
			if ptr, ok := args[0].(*int64); ok {
				*ptr = 1
			}
			// SupplierID
			if ptr, ok := args[1].(**int64); ok {
				supplierID := int64(10)
				*ptr = &supplierID
			}
			// ProductName
			if ptr, ok := args[2].(*string); ok {
				*ptr = "Produto Teste"
			}
			// Manufacturer
			if ptr, ok := args[3].(*string); ok {
				*ptr = "Fabricante Teste"
			}
			// Description
			if ptr, ok := args[4].(*string); ok {
				*ptr = "Descrição do produto"
			}
			// CostPrice
			if ptr, ok := args[5].(*float64); ok {
				*ptr = 50.0
			}
			// SalePrice
			if ptr, ok := args[6].(*float64); ok {
				*ptr = 100.0
			}
			// StockQuantity
			if ptr, ok := args[7].(*int); ok {
				*ptr = 100
			}
			// MinStock
			if ptr, ok := args[8].(*int); ok {
				*ptr = 10
			}
			// MaxStock
			if ptr, ok := args[9].(**int); ok {
				maxStock := 500
				*ptr = &maxStock
			}
			// Barcode
			if ptr, ok := args[10].(*string); ok {
				*ptr = "1234567890123"
			}
			// Status
			if ptr, ok := args[11].(*bool); ok {
				*ptr = true
			}
			// Version
			if ptr, ok := args[12].(*int); ok {
				*ptr = 1
			}
			// AllowDiscount
			if ptr, ok := args[13].(*bool); ok {
				*ptr = true
			}
			// MinDiscountPercent
			if ptr, ok := args[14].(*float64); ok {
				*ptr = 0.0
			}
			// MaxDiscountPercent
			if ptr, ok := args[15].(*float64); ok {
				*ptr = 30.0
			}
			// CreatedAt
			if ptr, ok := args[16].(*time.Time); ok {
				*ptr = now
			}
			// UpdatedAt
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

	t.Run("apply filters with ILIKE correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			ProductName:  "Dell",
			Manufacturer: "Dell",
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return !strings.Contains(query, "%%") && // Não deve ter %% escapado
				strings.Contains(query, "product_name ILIKE $1") &&
				strings.Contains(query, "manufacturer ILIKE $2")
		}), mock.MatchedBy(func(args []interface{}) bool {
			return len(args) >= 2 &&
				args[0] == "%Dell%" &&
				args[1] == "%Dell%"
		})).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("apply filters price ranges correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

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

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			found := 0
			for _, arg := range args {
				switch v := arg.(type) {
				case float64:
					if v == 50.0 || v == 100.0 || v == 200.0 {
						found++
					}
				}
			}
			return found >= 3
		})).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("apply filters supplier_id, status and allow_discount correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
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

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("apply filters barcode and version correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
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

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("apply filters date ranges correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		now := time.Now()
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

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			found := 0
			for _, arg := range args {
				switch v := arg.(type) {
				case time.Time:
					if v.Equal(createdFrom) || v.Equal(createdTo) ||
						v.Equal(updatedFrom) || v.Equal(updatedTo) {
						found++
					}
				}
			}
			return found >= 3
		})).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("uses allowed sort field when SortBy is valid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "sale_price",
				SortOrder: "desc",
			},
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(strings.ToLower(query), "order by sale_price desc")
		}), mock.Anything).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("apply filters min and max stock quantity correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		minStock := 10
		maxStock := 100

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			MinStockQuantity: &minStock,
			MaxStockQuantity: &maxStock,
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(query, "stock_quantity >= $") &&
				strings.Contains(query, "stock_quantity <= $")
		}), mock.MatchedBy(func(args []interface{}) bool {
			hasMinStock := false
			hasMaxStock := false

			for _, arg := range args {
				switch v := arg.(type) {
				case int:
					if v == 10 {
						hasMinStock = true
					} else if v == 100 {
						hasMaxStock = true
					}
				}
			}
			return hasMinStock && hasMaxStock
		})).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("apply filters only min stock quantity", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		minStock := 5

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			MinStockQuantity: &minStock,
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(query, "stock_quantity >= $") &&
				!strings.Contains(query, "stock_quantity <=")
		}), mock.MatchedBy(func(args []interface{}) bool {
			hasMinStock := false

			for _, arg := range args {
				if v, ok := arg.(int); ok && v == 5 {
					hasMinStock = true
				}
			}
			return hasMinStock
		})).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("apply filters only max stock quantity", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		maxStock := 50

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			MaxStockQuantity: &maxStock,
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return !strings.Contains(query, "stock_quantity >=") &&
				strings.Contains(query, "stock_quantity <= $")
		}), mock.MatchedBy(func(args []interface{}) bool {
			hasMaxStock := false

			for _, arg := range args {
				if v, ok := arg.(int); ok && v == 50 {
					hasMaxStock = true
				}
			}
			return hasMaxStock
		})).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("no stock filters applied when both nil", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			MinStockQuantity: nil,
			MaxStockQuantity: nil,
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return !strings.Contains(query, "stock_quantity >=") &&
				!strings.Contains(query, "stock_quantity <=")
		}), mock.Anything).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("defaults to created_at when SortBy is invalid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "invalid_field",
				SortOrder: "asc",
			},
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(strings.ToLower(query), "order by created_at asc")
		}), mock.Anything).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("defaults sort order to asc when SortOrder is invalid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "product_name",
				SortOrder: "INVALID",
			},
		}

		mockDB.On("Query", ctx, mock.MatchedBy(func(query string) bool {
			return strings.Contains(strings.ToLower(query), "order by product_name asc")
		}), mock.Anything).Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("database connection failed")

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(nil, dbErr)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("failed to scan row")
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(scanErr).Once()
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productFilterRepo{db: mockDB}
		ctx := context.Background()

		rowsErr := errors.New("iteration error")
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		filter := &filter.ProductFilter{
			BaseFilter: baseFilter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
