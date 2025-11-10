package repo

import (
	"context"
	"fmt"

	ifaceTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsgDb "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
	"github.com/jackc/pgx/v5"
)

type supplierCategoryRelationTx struct {
	db repo.DBExecutor
}

func NewSupplierCategoryRelationTx(db repo.DBExecutor) ifaceTx.SupplierCategoryRelationTx {
	return &supplierCategoryRelationTx{db: db}
}

func (r *supplierCategoryRelationTx) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
	const query = `
		INSERT INTO supplier_category_relations (supplier_id, category_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING created_at;
	`

	err := tx.QueryRow(ctx, query, relation.SupplierID, relation.CategoryID).Scan(&relation.CreatedAt)
	if err != nil {
		switch {
		case errMsgDb.IsDuplicateKey(err):
			return nil, errMsg.ErrRelationExists
		case errMsgDb.IsForeignKeyViolation(err):
			return nil, errMsg.ErrDBInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return relation, nil
}
