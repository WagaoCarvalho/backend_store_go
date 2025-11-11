package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *productCategoryRelationRepo) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error) {
	const query = `
		SELECT product_id, category_id, created_at
		FROM product_category_relations
		WHERE product_id = $1;
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.ProductCategoryRelation
	for rows.Next() {
		var rel models.ProductCategoryRelation
		if err := rows.Scan(&rel.ProductID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return relations, nil
}

func (r *productCategoryRelationRepo) HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error) {
	const query = `
		SELECT 1
		FROM product_category_relations
		WHERE product_id = $1 AND category_id = $2
		LIMIT 1;
	`

	var exists int
	err := r.db.QueryRow(ctx, query, productID, categoryID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return true, nil
}
