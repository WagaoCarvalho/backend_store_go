package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductRepo_Create(t *testing.T) {
	t.Run("successfully create product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			SupplierID:         utils.Int64Ptr(1),
			ProductName:        "Test Product",
			Manufacturer:       "Test Manufacturer",
			Description:        "Test Description",
			CostPrice:          10.50,
			SalePrice:          15.99,
			StockQuantity:      100,
			Barcode:            utils.StrToPtr("1234567890123"),
			Status:             true,
			AllowDiscount:      true,
			MaxDiscountPercent: 10.0,
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithID{
			IDValue:   1,
			TimeValue: now,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey on FK violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			SupplierID: utils.Int64Ptr(999), // ID inexistente
		}

		fkErr := errMsgPg.NewForeignKeyViolation("products_supplier_id_fkey")
		mockRow := &mockDb.MockRow{
			Err: fkErr,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate on unique violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			SupplierID:         utils.Int64Ptr(1),
			ProductName:        "Test Product",
			Manufacturer:       "Test Manufacturer",
			Description:        "Test Description",
			CostPrice:          10.50,
			SalePrice:          15.99,
			StockQuantity:      100,
			Barcode:            utils.StrToPtr("1234567890123"), // Barcode duplicado
			Status:             true,
			AllowDiscount:      true,
			MaxDiscountPercent: 10.0,
		}

		// Mock do erro de violação única (ex: barcode duplicado)
		uniqueErr := errMsgPg.NewUniqueViolation("products_barcode_key")
		mockRow := &mockDb.MockRow{
			Err: uniqueErr,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when query fails with generic database error", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{}
		dbErr := errors.New("database connection failed")

		mockRow := &mockDb.MockRow{
			Err: dbErr,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when pgx error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{}
		pgxErr := &pgconn.PgError{
			Message: "connection timeout",
		}

		mockRow := &mockDb.MockRow{
			Err: pgxErr,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, pgxErr.Message)
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_Update(t *testing.T) {
	t.Run("successfully update product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			ID:                 int64(1),
			SupplierID:         utils.Int64Ptr(1),
			ProductName:        "Updated Product",
			Manufacturer:       "Updated Manufacturer",
			Description:        "Updated Description",
			CostPrice:          12.50,
			SalePrice:          18.99,
			StockQuantity:      150,
			MinStock:           10,
			MaxStock:           utils.IntPtr(1000),
			Barcode:            utils.StrToPtr("9876543210987"),
			Status:             true,
			AllowDiscount:      true,
			MinDiscountPercent: 0.0,
			MaxDiscountPercent: 15.0,
			Version:            2,
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithID{
			IDValue:   2, // Nova versão
			TimeValue: now,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, product)

		assert.NoError(t, err)
		assert.Equal(t, now, product.UpdatedAt)
		assert.Equal(t, 2, product.Version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			ID:      int64(999),
			Version: 1,
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, product)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when FK violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			ID:         int64(1),
			SupplierID: utils.Int64Ptr(999), // ID inexistente
			Version:    1,
		}

		fkErr := errMsgPg.NewForeignKeyViolation("products_supplier_id_fkey")
		mockRow := &mockDb.MockRow{Err: fkErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, product)

		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrConflict when unique violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			ID:      int64(1),
			Version: 1,
		}

		uniqueErr := errMsgPg.NewUniqueViolation("products_barcode_key")
		mockRow := &mockDb.MockRow{Err: uniqueErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, product)

		assert.ErrorIs(t, err, errMsg.ErrConflict)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData when check violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			ID:      int64(1),
			Version: 1,
		}

		checkErr := errMsgPg.NewCheckViolation("products_sale_price_check")
		mockRow := &mockDb.MockRow{Err: checkErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, product)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when generic database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			ID:      int64(1),
			Version: 1,
		}

		dbErr := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		err := repo.Update(ctx, product)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_Delete(t *testing.T) {
	t.Run("successfully delete product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		cmdTag := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, productID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		cmdTag := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		dbError := errors.New("database connection failed")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(pgconn.CommandTag{}, dbError)

		err := repo.Delete(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
