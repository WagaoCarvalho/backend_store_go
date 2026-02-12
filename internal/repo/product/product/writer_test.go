package repo

import (
	"context"
	"errors"
	"strings"
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

	t.Run("return ErrInvalidData on check constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			SupplierID:         utils.Int64Ptr(1),
			ProductName:        "Test Product",
			Manufacturer:       "Test Manufacturer",
			Description:        "Test Description",
			CostPrice:          10.50,
			SalePrice:          -5.0, // Valor inválido (viola CHECK sale_price >= 0)
			StockQuantity:      100,
			MinStock:           0,
			MaxStock:           nil,
			Barcode:            utils.StrToPtr("1234567890123"),
			Status:             true,
			AllowDiscount:      true,
			MinDiscountPercent: 0.0,
			MaxDiscountPercent: 10.0,
		}

		// Cria erro de violação de CHECK (código 23514)
		pgErr := &pgconn.PgError{
			Code:    "23514", // Código específico para CHECK violation
			Message: "new row for relation \"products\" violates check constraint \"products_sale_price_check\"",
		}

		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)

		// Verifica que não é outro tipo de erro
		assert.NotErrorIs(t, err, errMsg.ErrDuplicate)
		assert.NotErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData on check constraint violation - discount range", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			SupplierID:         utils.Int64Ptr(1),
			ProductName:        "Test Product",
			Manufacturer:       "Test Manufacturer",
			Description:        "Test Description",
			CostPrice:          10.0,
			SalePrice:          20.0,
			StockQuantity:      100,
			MinStock:           0,
			MaxStock:           nil,
			Barcode:            utils.StrToPtr("1234567890123"),
			Status:             true,
			AllowDiscount:      true,
			MinDiscountPercent: 30.0, // Maior que max_discount_percent (viola CHECK chk_discount_range)
			MaxDiscountPercent: 20.0,
		}

		pgErr := &pgconn.PgError{
			Code:    "23514",
			Message: "new row for relation \"products\" violates check constraint \"chk_discount_range\"",
		}

		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData on check constraint violation - sale price less than cost", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()

		product := &models.Product{
			SupplierID:         utils.Int64Ptr(1),
			ProductName:        "Test Product",
			Manufacturer:       "Test Manufacturer",
			Description:        "Test Description",
			CostPrice:          20.0,
			SalePrice:          15.0, // Menor que cost_price (viola constraint implícita)
			StockQuantity:      100,
			MinStock:           0,
			MaxStock:           nil,
			Barcode:            utils.StrToPtr("1234567890123"),
			Status:             true,
			AllowDiscount:      true,
			MinDiscountPercent: 0.0,
			MaxDiscountPercent: 10.0,
		}

		pgErr := &pgconn.PgError{
			Code:    "23514",
			Message: "new row for relation \"products\" violates check constraint",
		}

		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, product)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
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
			ID:                 int64(999), // ID não existente
			SupplierID:         utils.Int64Ptr(1),
			ProductName:        "Product",
			Manufacturer:       "Manufacturer",
			Description:        "Description",
			CostPrice:          10.0,
			SalePrice:          15.0,
			StockQuantity:      100,
			MinStock:           5,
			MaxStock:           utils.IntPtr(500),
			Barcode:            utils.StrToPtr("12345678"),
			Status:             true,
			AllowDiscount:      false,
			MinDiscountPercent: 0.0,
			MaxDiscountPercent: 0.0,
			Version:            1,
		}

		// Mock para query principal retorna no rows
		mockRowMain := &mockDb.MockRow{Err: pgx.ErrNoRows}
		// Mock para verificação de existência também retorna no rows
		mockRowCheck := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRowMain)
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{product.ID}).
			Return(mockRowCheck)

		err := repo.Update(ctx, product)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrOptimisticLock when version mismatch", func(t *testing.T) {
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
			Version:            1, // Versão desatualizada
		}

		// 1. Mock para query UPDATE retorna no rows (versão não corresponde)
		mockRowUpdate := &mockDb.MockRow{Err: pgx.ErrNoRows}

		// 2. Mock para verificação confirma que produto EXISTE
		mockRowCheck := &mockDb.MockRow{Value: 1}

		// Chamada 1: UPDATE com WHERE id=$15 AND version=$16
		mockDB.On("QueryRow", ctx,
			mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "UPDATE products") && strings.Contains(query, "version = version + 1")
			}),
			mock.AnythingOfType("[]interface {}")).
			Return(mockRowUpdate).Once()

		// Chamada 2: SELECT 1 FROM products WHERE id = $1 (verificação)
		mockDB.On("QueryRow", ctx,
			mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "SELECT 1 FROM products WHERE id =")
			}),
			[]interface{}{product.ID}).
			Return(mockRowCheck).Once()

		err := repo.Update(ctx, product)

		assert.ErrorIs(t, err, errMsg.NotFoundOrErrVersionConflict)
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
