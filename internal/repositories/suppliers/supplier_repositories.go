package repository

import (
	"context"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	"github.com/jackc/pgx/v5/pgxpool"
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
		return models.Supplier{}, fmt.Errorf("erro ao criar fornecedor: %w", err)
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

	if err != nil {
		return nil, fmt.Errorf("fornecedor não encontrado: %w", err)
	}

	return &supplier, nil
}

func (r *supplierRepository) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	query := `SELECT id, name, cnpj, cpf, contact_info, created_at, updated_at FROM suppliers ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fornecedores: %w", err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		var s models.Supplier
		err := rows.Scan(&s.ID, &s.Name, &s.CNPJ, &s.CPF, &s.ContactInfo, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados do fornecedor: %w", err)
		}
		suppliers = append(suppliers, &s)
	}

	return suppliers, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	query := `
		UPDATE suppliers
		SET name = $1, cnpj = $2, cpf = $3, contact_info = $4, updated_at = $5
		WHERE id = $6
	`

	cmd, err := r.db.Exec(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.ContactInfo,
		time.Now(),
		supplier.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar fornecedor: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("fornecedor não encontrado")
	}

	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM suppliers WHERE id = $1`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar fornecedor: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("fornecedor não encontrado")
	}

	return nil
}
