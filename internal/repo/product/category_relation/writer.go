package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *productCategoryRelationRepo) Create(
	ctx context.Context,
	relation *models.ProductCategoryRelation,
) (*models.ProductCategoryRelation, error) {

	const query = `
        INSERT INTO product_category_relations (product_id, category_id, created_at)
        VALUES ($1, $2, NOW())
        RETURNING created_at;
    `

	err := r.db.QueryRow(
		ctx,
		query,
		relation.ProductID,
		relation.CategoryID,
	).Scan(&relation.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, fmt.Errorf("relação já existe [product_id=%d, category_id=%d]: %w",
					relation.ProductID, relation.CategoryID, errMsg.ErrRelationExists)
			case "23503":
				return nil, fmt.Errorf("chave estrangeira inválida [product_id=%d, category_id=%d]: %w",
					relation.ProductID, relation.CategoryID, errMsg.ErrDBInvalidForeignKey)
			}
		}

		return nil, fmt.Errorf("%w [product_id=%d, category_id=%d]: %v",
			errMsg.ErrCreate, relation.ProductID, relation.CategoryID, err)
	}

	return relation, nil
}

func (r *productCategoryRelationRepo) Delete(
	ctx context.Context,
	productID, categoryID int64,
) error {

	const query = `
		DELETE FROM product_category_relations
		WHERE product_id = $1 AND category_id = $2;
	`

	result, err := r.db.Exec(ctx, query, productID, categoryID)
	if err != nil {
		return fmt.Errorf("%w [product_id=%d, category_id=%d]: %w",
			errMsg.ErrDelete, productID, categoryID, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("relação não encontrada [product_id=%d, category_id=%d]: %w",
			productID, categoryID, errMsg.ErrNotFound)
	}

	return nil
}

func (r *productCategoryRelationRepo) DeleteAll(
	ctx context.Context,
	productID int64,
) error {

	const query = `
		DELETE FROM product_category_relations
		WHERE product_id = $1;
	`

	_, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("%w [product_id=%d]: %w",
			errMsg.ErrDelete, productID, err)
	}

	return nil
}
