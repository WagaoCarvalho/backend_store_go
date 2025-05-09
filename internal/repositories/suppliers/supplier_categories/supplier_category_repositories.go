package repositories

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierCategoryRepository interface {
	Create(ctx context.Context, category *models.SupplierCategory) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error)
	GetAll(ctx context.Context) ([]*models.SupplierCategory, error)
	Update(ctx context.Context, category *models.SupplierCategory) error
	Delete(ctx context.Context, id int64) error
}

type supplierCategoryRepository struct {
	db *pgxpool.Pool
}

func NewSupplierCategoryRepository(db *pgxpool.Pool) SupplierCategoryRepository {
	return &supplierCategoryRepository{db: db}
}

func (r *supplierCategoryRepository) Create(ctx context.Context, category *models.SupplierCategory) (int64, error) {
	query := `
		INSERT INTO supplier_categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar categoria: %w", err)
	}

	return category.ID, nil
}

func (r *supplierCategoryRepository) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM supplier_categories
		WHERE id = $1
	`

	var category models.SupplierCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	return &category, nil
}

func (r *supplierCategoryRepository) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM supplier_categories
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}
	defer rows.Close()

	var categories []*models.SupplierCategory
	for rows.Next() {
		var cat models.SupplierCategory
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}
		categories = append(categories, &cat)
	}

	return categories, nil
}

func (r *supplierCategoryRepository) Update(ctx context.Context, category *models.SupplierCategory) error {
	query := `
		UPDATE supplier_categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`

	cmdTag, err := r.db.Exec(ctx, query, category.Name, category.Description, category.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("categoria não encontrada")
	}

	return nil
}

func (r *supplierCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM supplier_categories WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("categoria não encontrada")
	}

	return nil
}
