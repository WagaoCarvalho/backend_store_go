package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsgDb "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5"
)

type SupplierCategoryRelation interface {
	Create(ctx context.Context, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error)
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierID(ctx context.Context, supplierID int64) error
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type supplierCategoryRelation struct {
	db repo.DBExecutor
}

func NewSupplierCategoryRelation(db repo.DBExecutor) SupplierCategoryRelation {
	return &supplierCategoryRelation{db: db}
}

func (r *supplierCategoryRelation) Create(ctx context.Context, rel *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
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

func (r *supplierCategoryRelation) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
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

func (r *supplierCategoryRelation) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	const query = `
		SELECT 1 FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2
		LIMIT 1;
	`

	var exists int
	err := r.db.QueryRow(ctx, query, supplierID, categoryID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return true, nil
}

func (r *supplierCategoryRelation) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error) {
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

	var rels []*models.SupplierCategoryRelation
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

func (r *supplierCategoryRelation) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error) {
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

	var rels []*models.SupplierCategoryRelation
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

func (r *supplierCategoryRelation) Delete(ctx context.Context, supplierID, categoryID int64) error {
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

func (r *supplierCategoryRelation) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
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
