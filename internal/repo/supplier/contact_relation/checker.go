package repo

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *supplierContactRelationRepo) HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error) {
	const query = `
		SELECT 1
		FROM supplier_contact_relations
		WHERE supplier_id = $1 AND contact_id = $2
		LIMIT 1;
	`

	var dummy int
	err := r.db.QueryRow(ctx, query, supplierID, contactID).Scan(&dummy)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return true, nil
}
