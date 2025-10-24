package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error)
	GetByID(ctx context.Context, id int64) (*models.ProductCategory, error)
	GetAll(ctx context.Context) ([]*models.ProductCategory, error)
	Update(ctx context.Context, category *models.ProductCategory) error
	Delete(ctx context.Context, id int64) error
}

type productCategoryRepository struct {
	db *pgxpool.Pool
}

func NewProductCategoryRepository(db *pgxpool.Pool) ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func (r *productCategoryRepository) Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	const query = `
		INSERT INTO product_categories (name, description, created_at, updated_at)
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

func (r *productCategoryRepository) GetByID(ctx context.Context, id int64) (*models.ProductCategory, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM product_categories
		WHERE id = $1;
	`

	var category models.ProductCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &category, nil
}

func (r *productCategoryRepository) GetAll(ctx context.Context) ([]*models.ProductCategory, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM product_categories
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var categories []*models.ProductCategory
	for rows.Next() {
		category := new(models.ProductCategory)
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

func (r *productCategoryRepository) Update(ctx context.Context, category *models.ProductCategory) error {
	const query = `
		UPDATE product_categories
		SET name = $1,
		    description = $2,
		    updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at;
	`

	err := r.db.QueryRow(ctx, query, category.Name, category.Description, category.ID).
		Scan(&category.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *productCategoryRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM product_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
