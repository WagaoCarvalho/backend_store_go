package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductCategoryRelation interface {
	Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error)
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error)
	HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error)
	GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error)
	Delete(ctx context.Context, productID, categoryID int64) error
	DeleteAll(ctx context.Context, productID int64) error
}

type productCategoryRelation struct {
	db *pgxpool.Pool
}

func NewProductCategoryRelation(db *pgxpool.Pool) ProductCategoryRelation {
	return &productCategoryRelation{db: db}
}

func (r *productCategoryRelation) Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	const query = `
		INSERT INTO product_category_relations (product_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.ProductID, relation.CategoryID)
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

func (r *productCategoryRelation) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error) {
	const query = `
		SELECT product_id, category_id, created_at
		FROM product_category_relations
		WHERE product_id = $1;
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.ProductCategoryRelation
	for rows.Next() {
		var rel models.ProductCategoryRelation
		if err := rows.Scan(&rel.ProductID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return relations, nil
}

func (r *productCategoryRelation) HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error) {
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

func (r *productCategoryRelation) Delete(ctx context.Context, productID, categoryID int64) error {
	const query = `
		DELETE FROM product_category_relations 
		WHERE product_id = $1 AND category_id = $2;
	`

	result, err := r.db.Exec(ctx, query, productID, categoryID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *productCategoryRelation) DeleteAll(ctx context.Context, productID int64) error {
	const query = `
		DELETE FROM product_category_relations
		WHERE product_id = $1;
	`

	_, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
