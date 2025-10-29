package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *sale) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	const query = `
		INSERT INTO sales (
			client_id, user_id, sale_date, total_amount, total_discount,
			payment_type, status, notes, version, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1, NOW(), NOW())
		RETURNING id, version, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		sale.ClientID,
		sale.UserID,
		sale.SaleDate,
		sale.TotalAmount,
		sale.TotalDiscount,
		sale.PaymentType,
		sale.Status,
		sale.Notes,
	).Scan(&sale.ID, &sale.Version, &sale.CreatedAt, &sale.UpdatedAt)

	if err != nil {
		switch {
		case errMsgPg.IsForeignKeyViolation(err):
			return nil, errMsg.ErrDBInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return sale, nil
}

func (r *sale) Update(ctx context.Context, sale *models.Sale) error {
	const query = `
		UPDATE sales
		SET 
			client_id      = $1,
			user_id        = $2,
			sale_date      = $3,
			total_amount   = $4,
			total_discount = $5,
			payment_type   = $6,
			status         = $7,
			notes          = $8,
			version        = version + 1,
			updated_at     = NOW()
		WHERE id = $9 AND version = $10
		RETURNING updated_at, version;
	`

	err := r.db.QueryRow(ctx, query,
		sale.ClientID,
		sale.UserID,
		sale.SaleDate,
		sale.TotalAmount,
		sale.TotalDiscount,
		sale.PaymentType,
		sale.Status,
		sale.Notes,
		sale.ID,
		sale.Version,
	).Scan(&sale.UpdatedAt, &sale.Version)

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

func (r *sale) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM sales WHERE id = $1;`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if tag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
