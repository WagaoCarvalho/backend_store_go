package repo

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

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
