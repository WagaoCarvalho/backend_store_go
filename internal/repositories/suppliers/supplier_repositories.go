package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error)
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id int64) error
}

type supplierRepository struct {
	db *pgxpool.Pool
}

func NewSupplierRepository(db *pgxpool.Pool) SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	query := `
		INSERT INTO suppliers (name, cnpj, cpf, contact_info)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.ContactInfo,
	).Scan(&supplier.ID, &supplier.CreatedAt, &supplier.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupplierCreate, err)
	}

	return supplier, nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, contact_info, created_at, updated_at
		FROM suppliers
		WHERE id = $1
	`

	var supplier models.Supplier
	err := r.db.QueryRow(ctx, query, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
		&supplier.ContactInfo,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSupplierNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupplierRetrieve, err)
	}

	return &supplier, nil
}

func (r *supplierRepository) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	const query = `
		SELECT id, name, cnpj, cpf, contact_info, created_at, updated_at
		FROM suppliers
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupplierList, err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		var s models.Supplier
		if err := rows.Scan(&s.ID, &s.Name, &s.CNPJ, &s.CPF, &s.ContactInfo, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSupplierList, err)
		}
		suppliers = append(suppliers, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupplierList, err)
	}

	return suppliers, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	const query = `
		UPDATE suppliers
		SET
			name         = $1,
			cnpj         = $2,
			cpf          = $3,
			contact_info = $4,
			updated_at   = NOW(),
			version      = version + 1
		WHERE
			id      = $5 AND
			version = $6
		RETURNING
			id,
			name,
			cnpj,
			cpf,
			contact_info,
			created_at,
			updated_at,
			version;
	`

	err := r.db.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.ContactInfo,
		supplier.ID,
		supplier.Version,
	).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
		&supplier.ContactInfo,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&supplier.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", ErrSupplierUpdate, err)
	}

	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM suppliers WHERE id = $1`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSupplierDelete, err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrSupplierNotFound
	}

	return nil
}
