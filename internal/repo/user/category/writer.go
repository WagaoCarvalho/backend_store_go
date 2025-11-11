package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *userCategoryRepo) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
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

func (r *userCategoryRepo) Update(ctx context.Context, category *models.UserCategory) error {
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

func (r *userCategoryRepo) Delete(ctx context.Context, id int64) error {
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
