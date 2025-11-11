package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *supplierContactRelationRepo) Create(ctx context.Context, relation *models.SupplierContactRelation) (*models.SupplierContactRelation, error) {
	const query = `
		INSERT INTO supplier_contact_relations (supplier_id, contact_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.SupplierID, relation.ContactID)
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

func (r *supplierContactRelationRepo) Delete(ctx context.Context, supplierID, contactID int64) error {
	const query = `
		DELETE FROM supplier_contact_relations
		WHERE supplier_id = $1 AND contact_id = $2;
	`

	_, err := r.db.Exec(ctx, query, supplierID, contactID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (r *supplierContactRelationRepo) DeleteAll(ctx context.Context, supplierID int64) error {
	const query = `
		DELETE FROM supplier_contact_relations
		WHERE supplier_id = $1;
	`

	_, err := r.db.Exec(ctx, query, supplierID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
