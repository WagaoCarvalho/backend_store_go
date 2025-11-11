package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

func (r *userCategoryRelationRepo) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelation, error) {
	if userID <= 0 {
		return []*models.UserCategoryRelation{}, errMsg.ErrZeroID
	}

	const query = `
        SELECT user_id, category_id, created_at
        FROM user_category_relations
        WHERE user_id = $1;
    `

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return []*models.UserCategoryRelation{}, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.UserCategoryRelation
	for rows.Next() {
		var rel models.UserCategoryRelation
		if err := rows.Scan(&rel.UserID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			return []*models.UserCategoryRelation{}, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return []*models.UserCategoryRelation{}, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	// Garantir que nunca retorne nil
	if relations == nil {
		relations = []*models.UserCategoryRelation{}
	}

	return relations, nil
}

func (r *userCategoryRelationRepo) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	const query = `
		SELECT 1
		FROM user_category_relations
		WHERE user_id = $1 AND category_id = $2
		LIMIT 1;
	`

	var exists int
	err := r.db.QueryRow(ctx, query, userID, categoryID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return true, nil
}
