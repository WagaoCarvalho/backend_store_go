package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

func (r *productRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	const query = `
		INSERT INTO products (
			supplier_id, product_name, manufacturer,
			product_description, cost_price, sale_price,
			stock_quantity, barcode, status,
			allow_discount, max_discount_percent,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		product.SupplierID,
		product.ProductName,
		product.Manufacturer,
		product.Description,
		product.CostPrice,
		product.SalePrice,
		product.StockQuantity,
		product.Barcode,
		product.Status,
		product.AllowDiscount,
		product.MaxDiscountPercent,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrDBInvalidForeignKey
		}

		if ok, constraint := errMsgPg.IsUniqueViolation(err); ok {
			return nil, fmt.Errorf("%w: %s", errMsg.ErrDuplicate, constraint)
		}

		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return product, nil
}

func (r *productRepo) Update(ctx context.Context, product *models.Product) error {
	const query = `
		UPDATE products
		SET
			supplier_id = $1,
			product_name = $2,
			manufacturer = $3,
			product_description = $4,
			cost_price = $5,
			sale_price = $6,
			stock_quantity = $7,
			min_stock = $8,
			max_stock = $9,
			barcode = $10,
			status = $11,
			version = version + 1,
			allow_discount = $12,
			min_discount_percent = $13,
			max_discount_percent = $14,
			updated_at = NOW()
		WHERE id = $15 AND version = $16
		RETURNING updated_at, version;
	`

	err := r.db.QueryRow(ctx, query,
		product.SupplierID,
		product.ProductName,
		product.Manufacturer,
		product.Description,
		product.CostPrice,
		product.SalePrice,
		product.StockQuantity,
		product.MinStock,
		product.MaxStock,
		product.Barcode,
		product.Status,
		product.AllowDiscount,
		product.MinDiscountPercent,
		product.MaxDiscountPercent,
		product.ID,
		product.Version,
	).Scan(&product.UpdatedAt, &product.Version)

	if err != nil {
		// Unique
		if ok, _ := errMsgPg.IsUniqueViolation(err); ok {
			return errMsg.ErrConflict
		}

		// Foreign key
		if errMsgPg.IsForeignKeyViolation(err) {
			return errMsg.ErrDBInvalidForeignKey
		}

		// Check constraint
		if errMsgPg.IsCheckViolation(err) {
			return errMsg.ErrInvalidData
		}

		// Not found
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}

		// Default
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *productRepo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM products WHERE id = $1;`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if tag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
