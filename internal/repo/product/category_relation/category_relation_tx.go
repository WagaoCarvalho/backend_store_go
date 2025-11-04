package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *productCategoryRelation) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	const query = `
		INSERT INTO product_category_relations (product_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := tx.Exec(ctx, query, relation.ProductID, relation.CategoryID)
	if err != nil {
		switch {
		case errMsgPg.IsDuplicateKey(err):
			return nil, errMsg.ErrRelationExists
		case errMsgPg.IsForeignKeyViolation(err):
			return nil, errMsg.ErrDBInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return relation, nil
}
