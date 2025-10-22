package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

func (r *productRepository) EnableDiscount(ctx context.Context, id int64) error {
	const query = `
		UPDATE products
		SET allow_discount = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductEnableDiscount, err)
	}

	return nil
}

func (r *productRepository) DisableDiscount(ctx context.Context, id int64) error {
	const query = `
		UPDATE products
		SET allow_discount = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductDisableDiscount, err)
	}

	return nil
}

func (r *productRepository) ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error) {
	const query = `
		UPDATE products
		SET max_discount_percent = $2, updated_at = NOW(), version = version + 1
		WHERE id = $1 AND allow_discount = TRUE
		RETURNING id, product_name, sale_price, max_discount_percent, allow_discount, version, updated_at;
	`

	var p models.Product
	if err := scanProductDiscountRow(r.db.QueryRow(ctx, query, id, percent), &p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			const checkQuery = `SELECT 1 FROM products WHERE id = $1`
			var exists int
			errCheck := r.db.QueryRow(ctx, checkQuery, id).Scan(&exists)
			if errCheck != nil || exists == 0 {
				return nil, errMsg.ErrNotFound
			}
			return nil, errMsg.ErrProductDiscountNotAllowed
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrProductApplyDiscount, err)
	}

	return &p, nil
}
