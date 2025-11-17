package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItemSale_GetByID(t *testing.T) {
	t.Run("successfully get item sale by id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		itemID := int64(1)
		expectedTime := time.Now()

		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []any{
				itemID,      // id
				int64(10),   // sale_id
				int64(20),   // product_id
				5,           // quantity
				10.5,        // unit_price
				1.0,         // discount
				0.5,         // tax
				50.0,        // subtotal
				"desc test", // description
				expectedTime,
				expectedTime,
			},
		}

		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{itemID}).
			Return(mockRow)

		result, err := repo.GetByID(ctx, itemID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, itemID, result.ID)
		assert.Equal(t, int64(10), result.SaleID)
		assert.Equal(t, int64(20), result.ProductID)
		assert.Equal(t, 5, result.Quantity)
		assert.Equal(t, 10.5, result.UnitPrice)
		assert.Equal(t, 1.0, result.Discount)
		assert.Equal(t, 0.5, result.Tax)
		assert.Equal(t, 50.0, result.Subtotal)
		assert.Equal(t, "desc test", result.Description)
		assert.Equal(t, expectedTime, result.CreatedAt)
		assert.Equal(t, expectedTime, result.UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when item does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		itemID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}
		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{itemID}).
			Return(mockRow)

		result, err := repo.GetByID(ctx, itemID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, pgx.ErrNoRows.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		itemID := int64(1)
		dbError := errors.New("connection lost")

		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.
			On("QueryRow", ctx, mock.Anything, []any{itemID}).
			Return(mockRow)

		result, err := repo.GetByID(ctx, itemID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})
}

func TestItemSale_GetBySaleID(t *testing.T) {
	t.Run("successfully get items by sale id with pagination", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(10)
		limit := 10
		offset := 0
		expectedTime := time.Now()

		// Criando mock rows com a estrutura correta
		mockRows := &mockDb.MockRows{
			Rows: []*mockDb.MockRow{
				{
					Values: []any{
						int64(1),      // id
						saleID,        // sale_id
						int64(20),     // product_id
						5,             // quantity
						10.5,          // unit_price
						1.0,           // discount
						0.5,           // tax
						50.0,          // subtotal
						"desc test 1", // description
						expectedTime,
						expectedTime,
					},
				},
				{
					Values: []any{
						int64(2),      // id
						saleID,        // sale_id
						int64(21),     // product_id
						3,             // quantity
						15.0,          // unit_price
						0.0,           // discount
						1.0,           // tax
						46.0,          // subtotal
						"desc test 2", // description
						expectedTime,
						expectedTime,
					},
				},
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, []any{saleID, limit, offset}).
			Return(mockRows, nil)

		results, err := repo.GetBySaleID(ctx, saleID, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 2)

		// Verify first item
		assert.Equal(t, int64(1), results[0].ID)
		assert.Equal(t, saleID, results[0].SaleID)
		assert.Equal(t, int64(20), results[0].ProductID)
		assert.Equal(t, 5, results[0].Quantity)
		assert.Equal(t, 10.5, results[0].UnitPrice)
		assert.Equal(t, 1.0, results[0].Discount)
		assert.Equal(t, 0.5, results[0].Tax)
		assert.Equal(t, 50.0, results[0].Subtotal)
		assert.Equal(t, "desc test 1", results[0].Description)
		assert.Equal(t, expectedTime, results[0].CreatedAt)
		assert.Equal(t, expectedTime, results[0].UpdatedAt)

		// Verify second item
		assert.Equal(t, int64(2), results[1].ID)
		assert.Equal(t, saleID, results[1].SaleID)
		assert.Equal(t, int64(21), results[1].ProductID)
		assert.Equal(t, 3, results[1].Quantity)
		assert.Equal(t, 15.0, results[1].UnitPrice)
		assert.Equal(t, 0.0, results[1].Discount)
		assert.Equal(t, 1.0, results[1].Tax)
		assert.Equal(t, 46.0, results[1].Subtotal)
		assert.Equal(t, "desc test 2", results[1].Description)
		assert.Equal(t, expectedTime, results[1].CreatedAt)
		assert.Equal(t, expectedTime, results[1].UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database query error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(10)
		limit := 10
		offset := 0
		dbError := errors.New("connection lost")

		mockDB.
			On("Query", ctx, mock.Anything, []any{saleID, limit, offset}).
			Return(nil, dbError)

		results, err := repo.GetBySaleID(ctx, saleID, limit, offset)

		assert.Nil(t, results)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get items with different pagination parameters", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(10)
		limit := 5
		offset := 10
		expectedTime := time.Now()

		mockRows := &mockDb.MockRows{
			Rows: []*mockDb.MockRow{
				{
					Values: []any{
						int64(11),     // id
						saleID,        // sale_id
						int64(30),     // product_id
						2,             // quantity
						8.0,           // unit_price
						0.5,           // discount
						0.3,           // tax
						15.8,          // subtotal
						"desc test 3", // description
						expectedTime,
						expectedTime,
					},
				},
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, []any{saleID, limit, offset}).
			Return(mockRows, nil)

		results, err := repo.GetBySaleID(ctx, saleID, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 1)

		assert.Equal(t, int64(11), results[0].ID)
		assert.Equal(t, saleID, results[0].SaleID)
		assert.Equal(t, int64(30), results[0].ProductID)
		assert.Equal(t, 2, results[0].Quantity)
		assert.Equal(t, 8.0, results[0].UnitPrice)
		assert.Equal(t, 0.5, results[0].Discount)
		assert.Equal(t, 0.3, results[0].Tax)
		assert.Equal(t, 15.8, results[0].Subtotal)
		assert.Equal(t, "desc test 3", results[0].Description)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when row scanning error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		saleID := int64(10)
		limit := 10
		offset := 0
		scanError := errors.New("scan error")

		// Criando mock rows que vai falhar no scan usando o campo Err do MockRow
		mockRows := &mockDb.MockRows{
			Rows: []*mockDb.MockRow{
				{
					Values: []any{
						int64(1),      // id
						saleID,        // sale_id
						int64(20),     // product_id
						5,             // quantity
						10.5,          // unit_price
						1.0,           // discount
						0.5,           // tax
						50.0,          // subtotal
						"desc test 1", // description
						time.Now(),
						time.Now(),
					},
					Err: scanError, // Isso fará o Scan retornar erro
				},
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, []any{saleID, limit, offset}).
			Return(mockRows, nil)

		results, err := repo.GetBySaleID(ctx, saleID, limit, offset)

		assert.Nil(t, results)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, scanError.Error())

		mockDB.AssertExpectations(t)
	})
}

func TestItemSale_GetByProductID(t *testing.T) {
	t.Run("successfully get items by product id with pagination", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(20)
		limit := 10
		offset := 0
		expectedTime := time.Now()

		mockRows := &mockDb.MockRows{
			Rows: []*mockDb.MockRow{
				{
					Values: []any{
						int64(1),      // id
						int64(10),     // sale_id
						productID,     // product_id
						5,             // quantity
						10.5,          // unit_price
						1.0,           // discount
						0.5,           // tax
						50.0,          // subtotal
						"desc test 1", // description
						expectedTime,
						expectedTime,
					},
				},
				{
					Values: []any{
						int64(2),      // id
						int64(11),     // sale_id
						productID,     // product_id
						3,             // quantity
						15.0,          // unit_price
						0.0,           // discount
						1.0,           // tax
						46.0,          // subtotal
						"desc test 2", // description
						expectedTime,
						expectedTime,
					},
				},
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, []any{productID, limit, offset}).
			Return(mockRows, nil)

		results, err := repo.GetByProductID(ctx, productID, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 2)

		// Verify first item
		assert.Equal(t, int64(1), results[0].ID)
		assert.Equal(t, int64(10), results[0].SaleID)
		assert.Equal(t, productID, results[0].ProductID)
		assert.Equal(t, 5, results[0].Quantity)
		assert.Equal(t, 10.5, results[0].UnitPrice)
		assert.Equal(t, 1.0, results[0].Discount)
		assert.Equal(t, 0.5, results[0].Tax)
		assert.Equal(t, 50.0, results[0].Subtotal)
		assert.Equal(t, "desc test 1", results[0].Description)
		assert.Equal(t, expectedTime, results[0].CreatedAt)
		assert.Equal(t, expectedTime, results[0].UpdatedAt)

		// Verify second item
		assert.Equal(t, int64(2), results[1].ID)
		assert.Equal(t, int64(11), results[1].SaleID)
		assert.Equal(t, productID, results[1].ProductID)
		assert.Equal(t, 3, results[1].Quantity)
		assert.Equal(t, 15.0, results[1].UnitPrice)
		assert.Equal(t, 0.0, results[1].Discount)
		assert.Equal(t, 1.0, results[1].Tax)
		assert.Equal(t, 46.0, results[1].Subtotal)
		assert.Equal(t, "desc test 2", results[1].Description)
		assert.Equal(t, expectedTime, results[1].CreatedAt)
		assert.Equal(t, expectedTime, results[1].UpdatedAt)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when database query error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(20)
		limit := 10
		offset := 0
		dbError := errors.New("connection lost")

		mockDB.
			On("Query", ctx, mock.Anything, []any{productID, limit, offset}).
			Return(nil, dbError)

		results, err := repo.GetByProductID(ctx, productID, limit, offset)

		assert.Nil(t, results)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbError.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when row scanning error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(20)
		limit := 10
		offset := 0
		scanError := errors.New("scan error")

		mockRows := &mockDb.MockRows{
			Rows: []*mockDb.MockRow{
				{
					Values: []any{
						int64(1),  // id
						int64(10), // sale_id
						productID, // product_id
						// ... outros campos
					},
					Err: scanError, // Isso fará o Scan retornar erro
				},
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, []any{productID, limit, offset}).
			Return(mockRows, nil)

		results, err := repo.GetByProductID(ctx, productID, limit, offset)

		assert.Nil(t, results)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, scanError.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get items with different pagination parameters", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &itemSaleRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(20)
		limit := 5
		offset := 10
		expectedTime := time.Now()

		mockRows := &mockDb.MockRows{
			Rows: []*mockDb.MockRow{
				{
					Values: []any{
						int64(11),     // id
						int64(15),     // sale_id
						productID,     // product_id
						2,             // quantity
						8.0,           // unit_price
						0.5,           // discount
						0.3,           // tax
						15.8,          // subtotal
						"desc test 3", // description
						expectedTime,
						expectedTime,
					},
				},
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, []any{productID, limit, offset}).
			Return(mockRows, nil)

		results, err := repo.GetByProductID(ctx, productID, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 1)

		assert.Equal(t, int64(11), results[0].ID)
		assert.Equal(t, int64(15), results[0].SaleID)
		assert.Equal(t, productID, results[0].ProductID)
		assert.Equal(t, 2, results[0].Quantity)
		assert.Equal(t, 8.0, results[0].UnitPrice)
		assert.Equal(t, 0.5, results[0].Discount)
		assert.Equal(t, 0.3, results[0].Tax)
		assert.Equal(t, 15.8, results[0].Subtotal)
		assert.Equal(t, "desc test 3", results[0].Description)

		mockDB.AssertExpectations(t)
	})
}
