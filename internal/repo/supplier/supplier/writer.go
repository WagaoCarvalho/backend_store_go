package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *supplierRepo) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	const query = `
		INSERT INTO suppliers (name, cnpj, cpf, description, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.Description,
		supplier.Status,
	).Scan(&supplier.ID, &supplier.CreatedAt, &supplier.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return supplier, nil
}

func (r *supplierRepo) Update(ctx context.Context, supplier *models.Supplier) error {
	const query = `
	UPDATE suppliers
	SET
		name        = $1,
		cnpj        = $2,
		cpf         = $3,
		description = $4,
		status      = $5,
		updated_at  = NOW(),
		version     = version + 1
	WHERE id = $6 AND version = $7
`

	result, err := r.db.Exec(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.Description,
		supplier.Status,
		supplier.ID,
		supplier.Version,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrZeroVersion
	}

	return nil

}

func (r *supplierRepo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM suppliers WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
