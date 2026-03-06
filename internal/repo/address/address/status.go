package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *addressRepo) Disable(ctx context.Context, aid int64) error {
	return r.setActive(ctx, aid, false)
}

func (r *addressRepo) Enable(ctx context.Context, aid int64) error {
	return r.setActive(ctx, aid, true)
}

// =======================
// Internal helper
// =======================

func (r *addressRepo) setActive(ctx context.Context, aid int64, active bool) error {

	const query = `
		UPDATE addresses
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at;
	`

	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, aid, active).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}

		if active {
			return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}
