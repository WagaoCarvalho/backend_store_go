package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

type ProductCategoryRelationRepoTx struct{}

func NewProductCategoryRelationRepoTx() *ProductCategoryRelationRepoTx {
	return &ProductCategoryRelationRepoTx{}
}

func (r *ProductCategoryRelationRepoTx) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.ProductCategoryRelation) error {
	const query = `
		INSERT INTO product_category_relations (product_id, category_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING created_at;
	`

	err := tx.QueryRow(ctx, query, relation.ProductID, relation.CategoryID).Scan(&relation.CreatedAt)
	if err != nil {
		switch {
		case errMsgPg.IsDuplicateKey(err):
			return fmt.Errorf("relação já existe [product_id=%d, category_id=%d]: %w",
				relation.ProductID, relation.CategoryID, errMsg.ErrRelationExists)

		case errMsgPg.IsForeignKeyViolation(err):
			return fmt.Errorf("chave estrangeira inválida [product_id=%d, category_id=%d]: %w",
				relation.ProductID, relation.CategoryID, errMsg.ErrDBInvalidForeignKey)

		default:
			return fmt.Errorf("%w [product_id=%d, category_id=%d]: %v",
				errMsg.ErrCreate, relation.ProductID, relation.CategoryID, err)
		}
	}

	return nil
}
