package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	err_msg_db "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierCategoryRelationRepository interface {
	Create(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error)
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierId(ctx context.Context, supplierID int64) error
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type supplierCategoryRelationRepo struct {
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewSupplierCategoryRelationRepo(db *pgxpool.Pool, lg logger.LoggerAdapterInterface) SupplierCategoryRelationRepository {
	return &supplierCategoryRelationRepo{db: db, logger: lg}
}

func (r *supplierCategoryRelationRepo) Create(ctx context.Context, rel *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	const ref = "[supplierCategoryRelationRepository - Create] - "
	const query = `
		INSERT INTO supplier_category_relations (supplier_id, category_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING created_at;
	`

	err := r.db.QueryRow(ctx, query, rel.SupplierID, rel.CategoryID).Scan(&rel.CreatedAt)
	if err != nil {
		switch {
		case err_msg_db.IsDuplicateKey(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyHasExists, map[string]any{
				"supplier_id": rel.SupplierID,
				"category_id": rel.CategoryID,
			})
			return nil, err_msg.ErrRelationExists

		case err_msg_db.IsForeignKeyViolation(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": rel.SupplierID,
				"category_id": rel.CategoryID,
			})
			return nil, err_msg.ErrInvalidForeignKey

		default:
			r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"supplier_id": rel.SupplierID,
				"category_id": rel.CategoryID,
			})
			return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
		}
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": rel.SupplierID,
		"category_id": rel.CategoryID,
	})
	return rel, nil
}

func (r *supplierCategoryRelationRepo) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	const ref = "[supplierCategoryRelationRepository - CreateTx] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"supplier_id": relation.SupplierID,
		"category_id": relation.CategoryID,
	})

	const query = `
		INSERT INTO supplier_category_relations (supplier_id, category_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING created_at
	`

	err := tx.QueryRow(ctx, query, relation.SupplierID, relation.CategoryID).
		Scan(&relation.CreatedAt)

	if err != nil {
		switch {
		case err_msg_db.IsDuplicateKey(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyHasExists, map[string]any{
				"supplier_id": relation.SupplierID,
				"category_id": relation.CategoryID,
			})
			return nil, err_msg.ErrRelationExists

		case err_msg_db.IsForeignKeyViolation(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": relation.SupplierID,
				"category_id": relation.CategoryID,
			})
			return nil, err_msg.ErrInvalidForeignKey

		default:
			r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"supplier_id": relation.SupplierID,
				"category_id": relation.CategoryID,
			})
			return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
		}
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": relation.SupplierID,
		"category_id": relation.CategoryID,
	})

	return relation, nil
}

func (r *supplierCategoryRelationRepo) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	const ref = "[supplierCategoryRelationRepository - HasRelation] - "
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
		r.logger.Error(ctx, err, ref+logger.LogVerificationError, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	r.logger.Info(ctx, ref+logger.LogVerificationSuccess, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})
	return true, nil
}

func (r *supplierCategoryRelationRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	const ref = "[supplierCategoryRelationRepository - GetBySupplierID] - "
	const query = `
		SELECT supplier_id, category_id, created_at
		FROM supplier_category_relations
		WHERE supplier_id = $1
		ORDER BY category_id;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"supplier_id": supplierID})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}
	defer rows.Close()

	var rels []*models.SupplierCategoryRelations
	for rows.Next() {
		rel := new(models.SupplierCategoryRelations)
		if err := rows.Scan(&rel.SupplierID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogScanError, map[string]any{"supplier_id": supplierID})
			return nil, fmt.Errorf("%w: %v", err_msg.ErrScan, err)
		}
		rels = append(rels, rel)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, map[string]any{"supplier_id": supplierID})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrScan, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"supplier_id": supplierID, "count": len(rels)})
	return rels, nil
}

func (r *supplierCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	const ref = "[supplierCategoryRelationRepository - GetByCategoryID] - "
	const query = `
		SELECT supplier_id, category_id, created_at
		FROM supplier_category_relations
		WHERE category_id = $1
		ORDER BY supplier_id;
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"category_id": categoryID})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}
	defer rows.Close()

	var rels []*models.SupplierCategoryRelations
	for rows.Next() {
		rel := new(models.SupplierCategoryRelations)
		if err := rows.Scan(&rel.SupplierID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogScanError, map[string]any{"category_id": categoryID})
			return nil, fmt.Errorf("%w: %v", err_msg.ErrScan, err)
		}
		rels = append(rels, rel)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, map[string]any{"category_id": categoryID})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrScan, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"category_id": categoryID, "count": len(rels)})
	return rels, nil
}

func (r *supplierCategoryRelationRepo) Delete(ctx context.Context, supplierID, categoryID int64) error {
	const ref = "[supplierCategoryRelationRepository - Delete] - "
	const query = `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2;
	`

	cmd, err := r.db.Exec(ctx, query, supplierID, categoryID)
	if err != nil {
		switch {
		case err_msg_db.IsForeignKeyViolation(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": supplierID,
				"category_id": categoryID,
			})
			return err_msg.ErrInvalidForeignKey
		default:
			r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
				"supplier_id": supplierID,
				"category_id": categoryID,
			})
			return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
		}
	}

	if cmd.RowsAffected() == 0 {
		return err_msg.ErrNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})
	return nil
}

func (r *supplierCategoryRelationRepo) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	const ref = "[supplierCategoryRelationRepository - DeleteAllBySupplierId] - "
	const query = `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1;
	`

	_, err := r.db.Exec(ctx, query, supplierID)
	if err != nil {
		switch {
		case err_msg_db.IsForeignKeyViolation(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": supplierID,
			})
			return err_msg.ErrInvalidForeignKey
		default:
			r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
				"supplier_id": supplierID,
			})
			return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
		}
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{"supplier_id": supplierID})
	return nil
}
