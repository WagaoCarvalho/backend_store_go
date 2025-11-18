package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *supplierCategoryRelationRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error) {
	const query = `
		SELECT supplier_id, category_id, created_at
		FROM supplier_category_relations
		WHERE supplier_id = $1
		ORDER BY category_id;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	// Inicializa com slice vazia em vez de nil
	rels := make([]*models.SupplierCategoryRelation, 0)

	for rows.Next() {
		rel := new(models.SupplierCategoryRelation)
		if err := rows.Scan(&rel.SupplierID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		rels = append(rels, rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
	}

	return rels, nil
}

func (r *supplierCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error) {
	const query = `
		SELECT supplier_id, category_id, created_at
		FROM supplier_category_relations
		WHERE category_id = $1
		ORDER BY supplier_id;
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	// Inicializa com slice vazia em vez de nil
	rels := make([]*models.SupplierCategoryRelation, 0)

	for rows.Next() {
		rel := new(models.SupplierCategoryRelation)
		if err := rows.Scan(&rel.SupplierID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		rels = append(rels, rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
	}

	return rels, nil
}
