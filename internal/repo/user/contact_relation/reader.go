package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *userContactRelationRepo) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelation, error) {
	const query = `
		SELECT user_id, contact_id, created_at
		FROM user_contact_relations
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.UserContactRelation
	for rows.Next() {
		var rel models.UserContactRelation
		if err := rows.Scan(&rel.UserID, &rel.ContactID, &rel.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if relations == nil {
		relations = []*models.UserContactRelation{}
	}

	return relations, nil
}
