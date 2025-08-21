package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
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
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewSupplierCategoryRepository(db *pgxpool.Pool, logger logger.LoggerAdapterInterface) SupplierCategoryRepository {
	return &supplierCategoryRepository{db: db, logger: logger}
}

func (r *supplierCategoryRepository) Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error) {
	ref := "[supplierCategoryRepository - Create] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name":        category.Name,
		"description": category.Description,
	})

	const query = `
		INSERT INTO supplier_categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query, category.Name, category.Description).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name":        category.Name,
			"description": category.Description,
		})
		return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryCreate, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"category_id": category.ID,
		"name":        category.Name,
	})

	return category, nil
}

func (r *supplierCategoryRepository) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	ref := "[supplierCategoryRepository - GetByID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"category_id": id,
	})

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
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"category_id": id,
			})
			return nil, ErrSupplierCategoryNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"category_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetCategoryByID, err)
	}
	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"category_id": category.ID,
		"name":        category.Name,
	})

	return &category, nil
}

func (r *supplierCategoryRepository) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	ref := "[supplierCategoryRepository - GetAll] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, nil)

	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM supplier_categories
		ORDER BY name;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryGetAll, err)
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
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, nil)
			return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryScanRow, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, nil)
		return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryIterate, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total": len(categories),
	})

	return categories, nil
}

func (r *supplierCategoryRepository) Update(ctx context.Context, category *models.SupplierCategory) error {
	ref := "[supplierCategoryRepository - Update] - "

	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"category_id": category.ID,
		"name":        category.Name,
	})

	const query = `
		UPDATE supplier_categories
		SET
			name        = $1,
			description = $2,
			updated_at  = NOW()
		WHERE
			id = $3
		RETURNING 
			updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Description,
		category.ID,
	).Scan(&category.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"category_id": category.ID,
			})
			return ErrSupplierCategoryNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"category_id": category.ID,
			"name":        category.Name,
		})
		return fmt.Errorf("%w: %v", ErrSupplierCategoryUpdate, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"category_id": category.ID,
	})

	return nil
}

func (r *supplierCategoryRepository) Delete(ctx context.Context, id int64) error {
	ref := "[supplierCategoryRepository - Delete] - "

	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"category_id": id,
	})

	const query = `DELETE FROM supplier_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"category_id": id,
		})
		return fmt.Errorf("%w: %v", ErrSupplierCategoryDelete, err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"category_id": id,
		})
		return ErrSupplierCategoryNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"category_id": id,
	})

	return nil
}
