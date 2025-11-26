package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *userContactRelationRepo) Create(ctx context.Context, relation *models.UserContactRelation) (*models.UserContactRelation, error) {
	const query = `
		INSERT INTO user_contact_relations (user_id, contact_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.UserID, relation.ContactID)
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

func (r *userContactRelationRepo) Delete(ctx context.Context, userID, contactID int64) error {
	const query = `
		DELETE FROM user_contact_relations
		WHERE user_id = $1 AND contact_id = $2;
	`

	result, err := r.db.Exec(ctx, query, userID, contactID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrIDNotFound
	}

	return nil
}

func (r *userContactRelationRepo) DeleteAll(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM user_contact_relations
		WHERE user_id = $1;
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
