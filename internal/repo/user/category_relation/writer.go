package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *userCategoryRelationRepo) Create(ctx context.Context, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error) {
	const query = `
		INSERT INTO user_category_relations (user_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.UserID, relation.CategoryID)
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

func (r *userCategoryRelationRepo) Delete(ctx context.Context, userID, categoryID int64) error {
	const query = `
		DELETE FROM user_category_relations 
		WHERE user_id = $1 AND category_id = $2;
	`

	result, err := r.db.Exec(ctx, query, userID, categoryID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *userCategoryRelationRepo) DeleteAll(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM user_category_relations
		WHERE user_id = $1;
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
