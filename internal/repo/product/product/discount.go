package repo

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

func (r *productRepo) EnableDiscount(ctx context.Context, id int64) error {
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

func (r *productRepo) DisableDiscount(ctx context.Context, id int64) error {
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

func (r *productRepo) ApplyDiscount(ctx context.Context, id int64, percent float64) error {
	// Validação básica do percentual
	if percent < 0 || percent > 100 {
		return errMsg.ErrInvalidDiscountPercent
	}

	const query = `
		UPDATE products
		SET max_discount_percent = $2, 
			min_discount_percent = LEAST($2, min_discount_percent),
			updated_at = NOW(), 
			version = version + 1
		WHERE id = $1 AND allow_discount = TRUE
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, percent).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Verifica se o produto existe
			const checkQuery = `SELECT 1 FROM products WHERE id = $1`
			var exists int
			if errCheck := r.db.QueryRow(ctx, checkQuery, id).Scan(&exists); errCheck != nil {
				return errMsg.ErrNotFound
			}
			// Produto existe mas não permite desconto
			return errMsg.ErrProductDiscountNotAllowed
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductApplyDiscount, err)
	}

	return nil
}
