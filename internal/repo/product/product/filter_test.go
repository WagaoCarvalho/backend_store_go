package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductRepo_Filter_AllFilters(t *testing.T) {
	t.Run("should return filtered products successfully with all filter types", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		// Mock filter data with various filter types
		minStock := 10
		maxStock := 100
		minCostPrice := 5.0
		maxCostPrice := 50.0
		status := true
		supplierID := int64(1)
		createdFrom := time.Now().AddDate(0, -1, 0)
		createdTo := time.Now()

		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  20,
				Offset: 0,
			},
			ProductName:      "Notebook",
			Manufacturer:     "Dell",
			Status:           &status,
			SupplierID:       &supplierID,
			MinCostPrice:     &minCostPrice,
			MaxCostPrice:     &maxCostPrice,
			MinStockQuantity: &minStock,
			MaxStockQuantity: &maxStock,
			CreatedFrom:      &createdFrom,
			CreatedTo:        &createdTo,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		// Mock 2 products being returned
		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Times(2).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestProductRepo_Filter_TextFilters(t *testing.T) {
	t.Run("should apply barcode filter when barcode is provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			Barcode: "1234567890123",
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestProductRepo_Filter_EqualFilters(t *testing.T) {
	t.Run("should apply version filter when version is provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		version := 2
		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			Version: &version,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("should apply allow_discount filter when allow_discount is provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		allowDiscount := true
		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			AllowDiscount: &allowDiscount,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestProductRepo_Filter_RangeFilters(t *testing.T) {
	t.Run("should apply sale_price range filter when min or max sale price is provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		minSalePrice := 50.0
		maxSalePrice := 200.0
		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			MinSalePrice: &minSalePrice,
			MaxSalePrice: &maxSalePrice,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("should apply discount percent range filter when min or max discount percent is provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		minDiscount := 5.0
		maxDiscount := 20.0
		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			MinDiscountPercent: &minDiscount,
			MaxDiscountPercent: &maxDiscount,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestProductRepo_Filter_DateFilters(t *testing.T) {
	t.Run("should apply updated_at range filter when both from and to dates are provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		updatedFrom := time.Now().AddDate(0, -1, 0)
		updatedTo := time.Now()
		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("should apply updated_at range filter when only from date is provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		updatedFrom := time.Now().AddDate(0, -1, 0)
		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   nil,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("should apply updated_at range filter when only to date is provided", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		updatedTo := time.Now()
		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			UpdatedFrom: nil,
			UpdatedTo:   &updatedTo,
		}

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.NotNil(t, products)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestProductRepo_Filter_ErrorScenarios(t *testing.T) {
	t.Run("should return error when database query fails", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			ProductName: "Test Product",
		}

		expectedErr := errors.New("database connection failed")

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(nil, expectedErr)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "erro ao buscar")
		assert.Contains(t, err.Error(), "database connection failed")

		mockDB.AssertExpectations(t)
	})

	t.Run("should return error when row scan fails", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			ProductName: "Test Product",
		}

		scanErr := errors.New("scan error: invalid data type")

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(scanErr)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "erro ao ler banco de dados")
		assert.Contains(t, err.Error(), "scan error: invalid data type")

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("should return error when rows iteration fails", func(t *testing.T) {
		// Setup
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)

		ctx := context.Background()
		repo := &productRepo{db: mockDB}

		filterData := &model.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			ProductName: "Test Product",
		}

		iterationErr := errors.New("iteration error: connection lost")

		// Mock expectations
		mockDB.On("Query", ctx, mock.AnythingOfType("string"), mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Once().Return(true)
		mockRows.On("Next").Once().Return(false)
		mockRows.On("Err").Return(iterationErr)
		mockRows.On("Close").Return()

		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("*int"),       // MinStock
			mock.AnythingOfType("**int"),      // MaxStock
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*int"),       // Version
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MinDiscountPercent
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		// Execute
		products, err := repo.Filter(ctx, filterData)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "erro ao iterar")
		assert.Contains(t, err.Error(), "iteration error: connection lost")

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
