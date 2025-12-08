package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

func (r *productRepo) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	const query = `
	SELECT id,
	       supplier_id,
	       product_name,
	       manufacturer,
	       product_description,
	       cost_price,
	       sale_price,
	       stock_quantity,
	       min_stock,
	       max_stock,
	       barcode,
	       status,
	       version,
	       allow_discount,
	       min_discount_percent,
	       max_discount_percent,
	       created_at,
	       updated_at
	FROM products
	WHERE id = $1;
	`

	var p models.Product
	if err := scanProductRow(r.db.QueryRow(ctx, query, id), &p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return &p, nil
}

func scanProductRow(row pgx.Row, p *models.Product) error {
	return row.Scan(
		&p.ID,
		&p.SupplierID,
		&p.ProductName,
		&p.Manufacturer,
		&p.Description,
		&p.CostPrice,
		&p.SalePrice,
		&p.StockQuantity,
		&p.MinStock,
		&p.MaxStock,
		&p.Barcode,
		&p.Status,
		&p.Version,
		&p.AllowDiscount,
		&p.MinDiscountPercent,
		&p.MaxDiscountPercent,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}

func scanProductRowLimited(row pgx.Row, p *models.Product) error {
	return row.Scan(
		&p.ID,
		&p.SupplierID,
		&p.ProductName,
		&p.Manufacturer,
		&p.Description,
		&p.CostPrice,
		&p.SalePrice,
		&p.StockQuantity,
		&p.Barcode,
		&p.Status,
		&p.AllowDiscount,
		&p.MaxDiscountPercent,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}

func scanProductDiscountRow(row pgx.Row, p *models.Product) error {
	return row.Scan(
		&p.ID,
		&p.ProductName,
		&p.SalePrice,
		&p.MaxDiscountPercent,
		&p.AllowDiscount,
		&p.Version,
		&p.UpdatedAt,
	)
}
