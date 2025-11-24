package repo

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *userContactRelationRepo) HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error) {
	const query = `
		SELECT 1
		FROM user_contact_relations
		WHERE user_id = $1 AND contact_id = $2
		LIMIT 1;
	`

	var dummy int
	err := r.db.QueryRow(ctx, query, userID, contactID).Scan(&dummy)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return true, nil
}
