package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5"
)

type UserCategory interface {
	Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	GetByID(ctx context.Context, id int64) (*models.UserCategory, error)
	GetAll(ctx context.Context) ([]*models.UserCategory, error)
	Update(ctx context.Context, category *models.UserCategory) error
	Delete(ctx context.Context, id int64) error
}

type userCategory struct {
	db repo.DBExecutor
}

func NewUserCategory(db repo.DBExecutor) UserCategory {
	return &userCategory{db: db}
}

func (r *userCategory) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	const query = `
		INSERT INTO user_categories (name, description, created_at, updated_at)
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

func (r *userCategory) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
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
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &category, nil
}

func (r *userCategory) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM user_categories
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var categories []*models.UserCategory
	for rows.Next() {
		category := new(models.UserCategory)
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

func (r *userCategory) Update(ctx context.Context, category *models.UserCategory) error {
	const query = `
		UPDATE user_categories
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

func (r *userCategory) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM user_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
