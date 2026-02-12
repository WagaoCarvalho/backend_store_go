package repo

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *productRepo) GetStock(ctx context.Context, id int64) (int, error) {
	const query = `
		SELECT stock_quantity
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

func (r *productRepo) UpdateStock(ctx context.Context, id int64, quantity int) error {
	if quantity < 0 {
		return errMsg.ErrInvalidQuantity
	}

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

func (r *productRepo) IncreaseStock(ctx context.Context, id int64, amount int) error {
	if amount <= 0 {
		return errMsg.ErrInvalidQuantity
	}

	const query = `
		UPDATE products
		SET stock_quantity = stock_quantity + $2, 
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

func (r *productRepo) DecreaseStock(ctx context.Context, id int64, amount int) error {
	if amount <= 0 {
		return errMsg.ErrInvalidQuantity
	}

	const query = `
		UPDATE products
		SET stock_quantity = stock_quantity - $2, 
		    updated_at = NOW(), 
		    version = version + 1
		WHERE id = $1 
		  AND stock_quantity >= $2  -- Garante que há estoque suficiente
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, amount).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Pode ser produto não encontrado OU estoque insuficiente
			// Verifica se o produto existe
			const checkQuery = `SELECT 1 FROM products WHERE id = $1`
			var exists int
			if errCheck := r.db.QueryRow(ctx, checkQuery, id).Scan(&exists); errCheck != nil {
				return errMsg.ErrNotFound
			}
			return errMsg.ErrInsufficientStock
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}
