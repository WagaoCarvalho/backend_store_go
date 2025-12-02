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

func (r *saleRepo) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	const query = `
		INSERT INTO sales (
			client_id,
			user_id,
			sale_date,
			total_items_amount,
			total_items_discount,
			total_sale_discount,
			total_amount,
			payment_type,
			status,
			notes,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		sale.ClientID,
		sale.UserID,
		sale.SaleDate,
		sale.TotalItemsAmount,
		sale.TotalItemsDiscount,
		sale.TotalSaleDiscount,
		sale.TotalAmount,
		sale.PaymentType,
		sale.Status,
		sale.Notes,
	).Scan(&sale.ID, &sale.CreatedAt, &sale.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrDBInvalidForeignKey
		}

		if ok, constraint := errMsgPg.IsUniqueViolation(err); ok {
			return nil, fmt.Errorf("%w: %s", errMsg.ErrDuplicate, constraint)
		}

		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return sale, nil
}

func (r *saleRepo) Update(ctx context.Context, sale *models.Sale) error {

	// 1) Seleciona versão atual
	const querySelect = `
		SELECT version
		FROM sales
		WHERE id = $1
	`

	var currentVersion int
	err := r.db.QueryRow(ctx, querySelect, sale.ID).Scan(&currentVersion)

	if errors.Is(err, pgx.ErrNoRows) {
		return errMsg.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: erro ao consultar venda: %v", errMsg.ErrUpdate, err)
	}

	// 2) Confere version (optimistic lock)
	if currentVersion != sale.Version {
		return errMsg.ErrZeroVersion
	}

	// 3) Atualiza
	const queryUpdate = `
		UPDATE sales
		SET 
			client_id            = $1,
			user_id              = $2,
			sale_date            = $3,
			total_items_amount   = $4,
			total_items_discount = $5,
			total_sale_discount  = $6,
			total_amount         = $7,
			payment_type         = $8,
			status               = $9,
			notes                = $10,
			updated_at           = NOW(),
			version              = version + 1
		WHERE id = $11
		RETURNING updated_at, version
	`

	err = r.db.QueryRow(ctx, queryUpdate,
		sale.ClientID,
		sale.UserID,
		sale.SaleDate,
		sale.TotalItemsAmount,
		sale.TotalItemsDiscount,
		sale.TotalSaleDiscount,
		sale.TotalAmount,
		sale.PaymentType,
		sale.Status,
		sale.Notes,
		sale.ID,
	).Scan(&sale.UpdatedAt, &sale.Version)

	if err != nil {
		switch {
		case errMsgPg.IsForeignKeyViolation(err):
			return errMsg.ErrDBInvalidForeignKey
		default:
			return fmt.Errorf("%w: erro ao atualizar venda: %v", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (r *saleRepo) Delete(ctx context.Context, id int64) error {
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
