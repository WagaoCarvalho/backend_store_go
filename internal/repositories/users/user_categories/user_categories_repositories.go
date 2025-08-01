package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserCategoryRepository interface {
	Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	GetByID(ctx context.Context, id int64) (*models.UserCategory, error)
	GetAll(ctx context.Context) ([]*models.UserCategory, error)
	Update(ctx context.Context, category *models.UserCategory) error
	Delete(ctx context.Context, id int64) error
}

type userCategoryRepository struct {
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewUserCategoryRepository(db *pgxpool.Pool, logger logger.LoggerAdapterInterface) UserCategoryRepository {
	return &userCategoryRepository{db: db, logger: logger}
}

func (r *userCategoryRepository) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	ref := "[userCategoryRepository - Create] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name":        category.Name,
		"description": category.Description,
	})

	const query = `
		INSERT INTO user_categories (name, description, created_at, updated_at)
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
		return nil, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"category_id": category.ID,
		"name":        category.Name,
	})

	return category, nil
}

func (r *userCategoryRepository) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
	ref := "[userCategoryRepository - GetByID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"category_id": id,
	})

	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM user_categories
		WHERE id = $1;
	`

	var category models.UserCategory

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
			return nil, ErrCategoryNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"category_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetCategoryByID, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"category_id": category.ID,
	})

	return &category, nil
}

func (r *userCategoryRepository) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	ref := "[userCategoryRepository - GetAll] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, nil)

	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM user_categories
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrGetCategories, err)
	}
	defer rows.Close()

	var categories []*models.UserCategory
	for rows.Next() {
		category := new(models.UserCategory)
		if err := rows.Scan(&category.ID, &category.Name, &category.Description,
			&category.CreatedAt, &category.UpdatedAt); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, nil)
			return nil, fmt.Errorf("%w: %v", ErrScanCategory, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, nil)
		return nil, fmt.Errorf("%w: %v", ErrIterateCategories, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total": len(categories),
	})

	return categories, nil
}

func (r *userCategoryRepository) Update(ctx context.Context, category *models.UserCategory) error {
	ref := "[userCategoryRepository - Update] - "

	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"category_id": category.ID,
		"name":        category.Name,
	})

	const query = `
		UPDATE user_categories
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
			return ErrCategoryNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"category_id": category.ID,
			"name":        category.Name,
		})
		return fmt.Errorf("%w: %v", ErrUpdateCategory, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"category_id": category.ID,
	})

	return nil
}

func (r *userCategoryRepository) Delete(ctx context.Context, id int64) error {
	ref := "[userCategoryRepository - Delete] - "
	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"category_id": id,
	})

	query := `DELETE FROM user_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"category_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"category_id": id,
		})
		return ErrCategoryNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"category_id": id,
	})

	return nil
}
