package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_contact_relations"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierContactRelationRepository interface {
	Create(ctx context.Context, relation *models.SupplierContactRelations) (*models.SupplierContactRelations, error)
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierContactRelations) (*models.SupplierContactRelations, error)
	HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error)
	GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelations, error)
	Delete(ctx context.Context, supplierID, contactID int64) error
	DeleteAll(ctx context.Context, supplierID int64) error
}

type supplierContactRelationRepositories struct {
	db *pgxpool.Pool
}

func NewSupplierContactRelationRepositories(db *pgxpool.Pool) SupplierContactRelationRepository {
	return &supplierContactRelationRepositories{db: db}
}

func (r *supplierContactRelationRepositories) Create(ctx context.Context, relation *models.SupplierContactRelations) (*models.SupplierContactRelations, error) {
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
			return nil, errMsg.ErrInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return relation, nil
}

func (r *supplierContactRelationRepositories) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierContactRelations) (*models.SupplierContactRelations, error) {
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
			return nil, errMsg.ErrInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return relation, nil
}

func (r *supplierContactRelationRepositories) HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error) {
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

func (r *supplierContactRelationRepositories) GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelations, error) {
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

	var relations []*models.SupplierContactRelations
	for rows.Next() {
		var rel models.SupplierContactRelations
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

func (r *supplierContactRelationRepositories) Delete(ctx context.Context, supplierID, contactID int64) error {
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

func (r *supplierContactRelationRepositories) DeleteAll(ctx context.Context, supplierID int64) error {
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
