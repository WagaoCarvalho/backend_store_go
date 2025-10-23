package repo

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *product) UpdateStock(ctx context.Context, id int64, quantity int) error {
	const query = `
		UPDATE products
		SET stock_quantity = $2, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, quantity).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *product) IncreaseStock(ctx context.Context, id int64, amount int) error {
	const query = `
		UPDATE products
		SET stock_quantity = stock_quantity + $2, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, amount).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *product) DecreaseStock(ctx context.Context, id int64, amount int) error {
	const query = `
		UPDATE products
		SET stock_quantity = GREATEST(COALESCE(stock_quantity, 0) - $2, 0),
		    updated_at = NOW(),
		    version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, amount).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *product) GetStock(ctx context.Context, id int64) (int, error) {
	const query = `
		SELECT COALESCE(stock_quantity, 0)
		FROM products
		WHERE id = $1;
	`

	var stock int
	err := r.db.QueryRow(ctx, query, id).Scan(&stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errMsg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return stock, nil
}
