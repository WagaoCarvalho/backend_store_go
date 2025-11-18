package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *supplierContactRelationRepo) GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelation, error) {
	const query = `
		SELECT supplier_id, contact_id, created_at
		FROM supplier_contact_relations
		WHERE supplier_id = $1;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.SupplierContactRelation
	for rows.Next() {
		var rel models.SupplierContactRelation
		if err := rows.Scan(&rel.SupplierID, &rel.ContactID, &rel.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return relations, nil
}
