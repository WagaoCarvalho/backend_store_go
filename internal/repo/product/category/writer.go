package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *productCategoryRepo) Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error) {
	const query = `
		INSERT INTO product_categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query, category.Name, category.Description).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if ok, constraint := errMsgPg.IsUniqueViolation(err); ok {
			return nil, fmt.Errorf("%w: %s", errMsg.ErrDuplicate, constraint)
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return category, nil
}

// No método Update, adicionar verificação de conflito de nome único
func (r *productCategoryRepo) Update(ctx context.Context, category *models.ProductCategory) error {
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
		// Verificar se é violação de unique (exceto para o próprio registro)
		if ok, constraint := errMsgPg.IsUniqueViolation(err); ok {
			return fmt.Errorf("%w: %s", errMsg.ErrDuplicate, constraint)
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *productCategoryRepo) Delete(ctx context.Context, id int64) error {
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
