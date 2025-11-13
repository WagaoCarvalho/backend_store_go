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

func (r *productRepo) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	const query = `SELECT version FROM products WHERE id = $1`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errMsg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGetVersion, err)
	}

	return version, nil
}
