package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

func (r *productRepo) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	const query = `
	SELECT id, supplier_id, product_name, manufacturer, product_description,
		cost_price, sale_price, stock_quantity, barcode,
		status, allow_discount, max_discount_percent,
		created_at, updated_at
	FROM products
	ORDER BY id
	LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProductRowLimited(rows, &p); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return products, nil
}

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

func (r *productRepo) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
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
	WHERE product_name ILIKE '%' || $1 || '%'
	ORDER BY product_name;
	`

	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProductRow(rows, &p); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return products, nil
}

func (r *productRepo) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {
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
	WHERE manufacturer ILIKE '%' || $1 || '%'
	ORDER BY product_name;
	`

	rows, err := r.db.Query(ctx, query, manufacturer)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProductRow(rows, &p); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return products, nil
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

func (r *productRepo) ProductExists(ctx context.Context, productID int64) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, productID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return exists, nil
}
