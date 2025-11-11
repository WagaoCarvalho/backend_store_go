package repo

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *saleRepo) Cancel(ctx context.Context, id int64) error {
	const query = `
		UPDATE sales
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = $1;
	`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	if tag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}
	return nil
}

func (r *saleRepo) Complete(ctx context.Context, id int64) error {
	const query = `
		UPDATE sales
		SET status = 'completed', updated_at = NOW()
		WHERE id = $1;
	`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	if tag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}
	return nil
}
