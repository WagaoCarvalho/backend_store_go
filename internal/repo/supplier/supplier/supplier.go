package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error)
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetByName(ctx context.Context, name string) ([]*models.Supplier, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
	SupplierExists(ctx context.Context, supplierID int64) (bool, error)
}

type supplierRepository struct {
	db *pgxpool.Pool
}

func NewSupplierRepository(db *pgxpool.Pool) SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	const query = `
		INSERT INTO suppliers (name, cnpj, cpf, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.Status,
	).Scan(&supplier.ID, &supplier.CreatedAt, &supplier.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return supplier, nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, status, created_at, updated_at
		FROM suppliers
		WHERE id = $1
	`

	var supplier models.Supplier
	err := r.db.QueryRow(ctx, query, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
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

func (r *supplierRepository) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, status, created_at, updated_at
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

func (r *supplierRepository) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, status, created_at, updated_at
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

func (r *supplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	const query = `
		UPDATE suppliers
		SET
			name       = $1,
			cnpj       = $2,
			cpf        = $3,
			status     = $4,
			updated_at = NOW(),
			version    = version + 1
		WHERE
			id      = $5 AND
			version = $6
		RETURNING id, name, cnpj, cpf, status, created_at, updated_at, version
	`

	err := r.db.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.Status,
		supplier.ID,
		supplier.Version,
	).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
		&supplier.Status,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&supplier.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id int64) error {
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

func (r *supplierRepository) Disable(ctx context.Context, id int64) error {
	const query = `
		UPDATE suppliers
		SET status = false,
		    updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *supplierRepository) Enable(ctx context.Context, id int64) error {
	const query = `
		UPDATE suppliers
		SET status = true,
		    updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *supplierRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
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

func (r *supplierRepository) SupplierExists(ctx context.Context, supplierID int64) (bool, error) {
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
