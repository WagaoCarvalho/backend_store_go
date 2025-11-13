package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductRepo_GetAll(t *testing.T) {
	t.Run("successfully get all products - empty result", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)
		mockRows.On("Next").Return(false) // No rows
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		products, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		// products pode ser nil ou slice vazia, dependendo da implementação
		if products == nil {
			assert.Nil(t, products) // Aceita nil
		} else {
			assert.Empty(t, products) // Ou slice vazia
		}
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get all products - with data", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)
		mockRows.On("Next").Once().Return(true)  // First row
		mockRows.On("Next").Once().Return(false) // No more rows
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
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(nil)

		products, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		queryErr := errors.New("database connection failed")

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, queryErr)

		products, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when row scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan error")

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)
		mockRows.On("Next").Return(true)
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // ID
			mock.AnythingOfType("**int64"),    // SupplierID
			mock.AnythingOfType("*string"),    // ProductName
			mock.AnythingOfType("*string"),    // Manufacturer
			mock.AnythingOfType("*string"),    // Description
			mock.AnythingOfType("*float64"),   // CostPrice
			mock.AnythingOfType("*float64"),   // SalePrice
			mock.AnythingOfType("*int"),       // StockQuantity
			mock.AnythingOfType("**string"),   // Barcode
			mock.AnythingOfType("*bool"),      // Status
			mock.AnythingOfType("*bool"),      // AllowDiscount
			mock.AnythingOfType("*float64"),   // MaxDiscountPercent
			mock.AnythingOfType("*time.Time"), // CreatedAt
			mock.AnythingOfType("*time.Time"), // UpdatedAt
		).Return(scanErr)
		mockRows.On("Close").Return()

		products, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return error when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		mockRows := new(mockDb.MockRows)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		iterationErr := errors.New("iteration error")

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(iterationErr)
		mockRows.On("Close").Return()

		products, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
