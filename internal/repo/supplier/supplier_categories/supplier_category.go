package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierCategoryRepository interface {
	Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error)
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

func (r *supplierCategoryRepository) Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error) {
	const query = `
		INSERT INTO supplier_categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query, category.Name, category.Description).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return category, nil
}

func (r *supplierCategoryRepository) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM supplier_categories
		WHERE id = $1;
	`

	var category models.SupplierCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errMsg.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &category, nil
}

func (r *supplierCategoryRepository) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM supplier_categories
		ORDER BY name;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var categories []*models.SupplierCategory
	for rows.Next() {
		category := new(models.SupplierCategory)
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return categories, nil
}

func (r *supplierCategoryRepository) Update(ctx context.Context, category *models.SupplierCategory) error {
	const query = `
		UPDATE supplier_categories
		SET
			name        = $1,
			description = $2,
			updated_at  = NOW()
		WHERE id = $3
		RETURNING updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Description,
		category.ID,
	).Scan(&category.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *supplierCategoryRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM supplier_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
