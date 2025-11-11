package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *supplierRepo) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, status, description, created_at, updated_at
		FROM suppliers
		WHERE id = $1
	`

	var supplier models.Supplier
	err := r.db.QueryRow(ctx, query, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
		&supplier.Description,
		&supplier.Status,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errMsg.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &supplier, nil
}

func (r *supplierRepo) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, description, status, created_at, updated_at
		FROM suppliers
		WHERE name ILIKE $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	suppliers := make([]*models.Supplier, 0, 10)
	for rows.Next() {
		s := new(models.Supplier)
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.CNPJ,
			&s.CPF,
			&s.Description,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		suppliers = append(suppliers, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(suppliers) == 0 {
		return nil, errMsg.ErrNotFound
	}

	return suppliers, nil
}

func (r *supplierRepo) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, description, status, created_at, updated_at
		FROM suppliers
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	suppliers := make([]*models.Supplier, 0, 10)
	for rows.Next() {
		var s models.Supplier
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.CNPJ,
			&s.CPF,
			&s.Description,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		suppliers = append(suppliers, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(suppliers) == 0 {
		return nil, errMsg.ErrNotFound
	}

	return suppliers, nil
}

func (r *supplierRepo) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	const query = `
		SELECT version
		FROM suppliers
		WHERE id = $1
	`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errMsg.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return version, nil
}

func (r *supplierRepo) SupplierExists(ctx context.Context, supplierID int64) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM suppliers
			WHERE id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, supplierID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return exists, nil
}
