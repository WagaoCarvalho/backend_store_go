package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSupplierCategoryNotFound = errors.New("categoria de fornecedor não encontrada")
	ErrSupplierCategoryCreate   = errors.New("erro ao criar categoria")
	ErrSupplierCategoryGetAll   = errors.New("erro ao buscar categorias")
	ErrSupplierCategoryScanRow  = errors.New("erro ao ler dados da categoria")
	ErrSupplierCategoryUpdate   = errors.New("erro ao atualizar categoria")
	ErrSupplierCategoryDelete   = errors.New("erro ao deletar categoria")
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
		return 0, fmt.Errorf("%w: %v", ErrSupplierCategoryCreate, err)
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
		return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryNotFound, err)
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
		return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryGetAll, err)
	}
	defer rows.Close()

	var categories []*models.SupplierCategory
	for rows.Next() {
		var cat models.SupplierCategory
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryScanRow, err)
		}
		categories = append(categories, &cat)
	}

	return categories, nil
}

func (r *supplierCategoryRepository) Update(ctx context.Context, category *models.SupplierCategory) error {
	const query = `
		UPDATE supplier_categories
		SET
			name        = $1,
			description = $2,
			updated_at  = NOW(),
			version     = version + 1
		WHERE
			id      	= $3 AND
			version 	= $4
		RETURNING
			id,
			name,
			description,
			created_at,
			updated_at,
			version;
	`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Description,
		category.ID,
		category.Version,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSupplierCategoryNotFound
		}
		return fmt.Errorf("%w: %v", ErrSupplierCategoryUpdate, err)
	}

	return nil
}

func (r *supplierCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM supplier_categories WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSupplierCategoryDelete, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrSupplierCategoryNotFound
	}

	return nil
}
