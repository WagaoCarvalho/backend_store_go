package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

const baseSelectSupplier = `
	SELECT 
		id, name, cnpj, cpf, description, version, status, created_at, updated_at
	FROM suppliers
`

func (r *supplierRepo) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	const query = baseSelectSupplier + ` WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	supplier, err := scanSupplier(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return supplier, nil
}

func (r *supplierRepo) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	const query = baseSelectSupplier + ` WHERE name ILIKE $1 ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		s, err := scanSupplier(rows)
		if err != nil {
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
	const query = baseSelectSupplier + ` ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		s, err := scanSupplier(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		suppliers = append(suppliers, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
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

// scanSupplier é a função auxiliar que segue o padrão do scanAddress
type scanner interface {
	Scan(dest ...any) error
}

func scanSupplier(s scanner) (*models.Supplier, error) {
	var supplier models.Supplier

	err := s.Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
		&supplier.Description,
		&supplier.Version,
		&supplier.Status,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &supplier, nil
}
