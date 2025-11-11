package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsgDb "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *supplierCategoryRelationRepo) Create(ctx context.Context, rel *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
	const query = `
		INSERT INTO supplier_category_relations (supplier_id, category_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING created_at;
	`

	err := r.db.QueryRow(ctx, query, rel.SupplierID, rel.CategoryID).Scan(&rel.CreatedAt)
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

	return rel, nil
}

func (r *supplierCategoryRelationRepo) Delete(ctx context.Context, supplierID, categoryID int64) error {
	const query = `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2;
	`

	cmd, err := r.db.Exec(ctx, query, supplierID, categoryID)
	if err != nil {
		if errMsgDb.IsForeignKeyViolation(err) {
			return errMsg.ErrDBInvalidForeignKey
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if cmd.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *supplierCategoryRelationRepo) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
	const query = `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1;
	`

	_, err := r.db.Exec(ctx, query, supplierID)
	if err != nil {
		if errMsgDb.IsForeignKeyViolation(err) {
			return errMsg.ErrDBInvalidForeignKey
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
