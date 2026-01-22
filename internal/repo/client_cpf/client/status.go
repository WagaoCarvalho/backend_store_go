package repo

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *clientCpfRepo) setStatus(
	ctx context.Context,
	id int64,
	status bool,
) error {
	const query = `
		UPDATE clients_cpf
		SET
			status = $1,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $2
		RETURNING id;
	`

	var returnedID int64

	err := r.db.QueryRow(ctx, query, status, id).Scan(&returnedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (r *clientCpfRepo) Disable(ctx context.Context, id int64) error {
	return r.setStatus(ctx, id, false)
}

func (r *clientCpfRepo) Enable(ctx context.Context, id int64) error {
	return r.setStatus(ctx, id, true)
}
