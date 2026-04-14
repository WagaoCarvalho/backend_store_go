package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *supplierRepo) Disable(ctx context.Context, id int64) error {
	return r.setStatus(ctx, id, false)
}

func (r *supplierRepo) Enable(ctx context.Context, id int64) error {
	return r.setStatus(ctx, id, true)
}

func (r *supplierRepo) setStatus(ctx context.Context, id int64, status bool) error {
	const query = `
		UPDATE suppliers
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id, status).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}

		if status {
			return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}
