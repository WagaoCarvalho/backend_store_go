package repo

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *clientRepo) Disable(ctx context.Context, id int64) error {
	const query = `UPDATE clients SET status=false, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}

func (r *clientRepo) Enable(ctx context.Context, id int64) error {
	const query = `UPDATE clients SET status=true, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}
