package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *supplierCategoryRepo) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
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

func (r *supplierCategoryRepo) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
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

	// Inicializa com slice vazia em vez de nil
	categories := make([]*models.SupplierCategory, 0)

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
