package repo

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *clientCpfRepo) ClientCpfExists(ctx context.Context, clientID int64) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM clients WHERE id=$1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, clientID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return exists, nil
}
