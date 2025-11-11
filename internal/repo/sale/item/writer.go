package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *saleItemRepo) Create(ctx context.Context, item *models.SaleItem) (*models.SaleItem, error) {
	const query = `
		INSERT INTO sale_items (
			sale_id, product_id, quantity, unit_price, discount, tax, subtotal, description, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		item.SaleID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice,
		item.Discount,
		item.Tax,
		item.Subtotal,
		item.Description,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		switch {
		case errMsgPg.IsForeignKeyViolation(err):
			return nil, errMsg.ErrDBInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return item, nil
}

func (r *saleItemRepo) Update(ctx context.Context, item *models.SaleItem) error {
	const query = `
		UPDATE sale_items
		SET 
			sale_id     = $1,
			product_id  = $2,
			quantity    = $3,
			unit_price  = $4,
			discount    = $5,
			tax         = $6,
			subtotal    = $7,
			description = $8,
			updated_at  = NOW()
		WHERE id = $9
		RETURNING updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		item.SaleID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice,
		item.Discount,
		item.Tax,
		item.Subtotal,
		item.Description,
		item.ID,
	).Scan(&item.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return errMsg.ErrNotFound
		case errMsgPg.IsForeignKeyViolation(err):
			return errMsg.ErrDBInvalidForeignKey
		default:
			return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (r *saleItemRepo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM sale_items WHERE id = $1;`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if tag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *saleItemRepo) DeleteBySaleID(ctx context.Context, saleID int64) error {
	const query = `DELETE FROM sale_items WHERE sale_id = $1;`

	_, err := r.db.Exec(ctx, query, saleID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
