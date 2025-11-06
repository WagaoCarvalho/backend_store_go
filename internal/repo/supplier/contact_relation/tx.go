package repo

import (
	"context"
	"fmt"

	ifaceTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5"
)

type supplierContactRelationTx struct {
	db repo.DBExecutor
}

func NewSupplierContactRelationTx(db repo.DBExecutor) ifaceTx.SupplierContactRelationTx {
	return &supplierContactRelationTx{db: db}
}

func (r *supplierContactRelationTx) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierContactRelation) (*models.SupplierContactRelation, error) {
	const query = `
		INSERT INTO supplier_contact_relations (supplier_id, contact_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := tx.Exec(ctx, query, relation.SupplierID, relation.ContactID)
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
