package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSupplierNotFound    = errors.New("fornecedor não encontrado")
	ErrSupplierCreate      = errors.New("erro ao criar fornecedor")
	ErrSupplierUpdate      = errors.New("erro ao atualizar fornecedor")
	ErrSupplierDelete      = errors.New("erro ao deletar fornecedor")
	ErrSupplierList        = errors.New("erro ao listar fornecedores")
	ErrSupplierRetrieve    = errors.New("erro ao buscar fornecedor por ID")
	ErrInvalidSupplierData = errors.New("dados inválidos para fornecedor")
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier models.Supplier) (models.Supplier, error)
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

func (r *supplierRepository) Create(ctx context.Context, supplier models.Supplier) (models.Supplier, error) {
	if supplier.Name == "" || (supplier.CNPJ == nil && supplier.CPF == nil) {
		return models.Supplier{}, ErrInvalidSupplierData
	}

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
		return models.Supplier{}, fmt.Errorf("%w: %v", ErrSupplierCreate, err)
	}

	return supplier, nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	query := `SELECT id, name, cnpj, cpf, contact_info, created_at, updated_at FROM suppliers WHERE id = $1`

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
	query := `SELECT id, name, cnpj, cpf, contact_info, created_at, updated_at FROM suppliers ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupplierList, err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		var s models.Supplier
		err := rows.Scan(&s.ID, &s.Name, &s.CNPJ, &s.CPF, &s.ContactInfo, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSupplierList, err)
		}
		suppliers = append(suppliers, &s)
	}

	return suppliers, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	if supplier.ID <= 0 || supplier.Name == "" {
		return ErrInvalidSupplierData
	}

	query := `
		UPDATE suppliers
		SET name = $1, cnpj = $2, cpf = $3, contact_info = $4, updated_at = $5
		WHERE id = $6
	`

	now := time.Now()

	cmd, err := r.db.Exec(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.ContactInfo,
		now,
		supplier.ID,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrSupplierUpdate, err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrSupplierNotFound
	}

	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM suppliers WHERE id = $1`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSupplierDelete, err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrSupplierNotFound
	}

	return nil
}
