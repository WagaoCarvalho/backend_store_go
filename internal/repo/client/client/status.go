package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *clientRepo) Disable(ctx context.Context, id int64) error {
	const query = `
		UPDATE clients
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (r *clientRepo) Enable(ctx context.Context, id int64) error {
	const query = `
		UPDATE clients
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}
