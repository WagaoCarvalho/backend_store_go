package repo

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *supplierRepo) Disable(ctx context.Context, id int64) error {
	const query = `
		UPDATE suppliers
		SET status = false,
		    updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *supplierRepo) Enable(ctx context.Context, id int64) error {
	const query = `
		UPDATE suppliers
		SET status = true,
		    updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
