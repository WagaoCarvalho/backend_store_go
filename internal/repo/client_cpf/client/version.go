package repo

import (
	"context"
	"errors"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *clientCpfRepo) GetVersionByID(ctx context.Context, id int64) (int, error) {
	const query = `
		SELECT version
		FROM clients_cpf
		WHERE id = $1
		LIMIT 1
	`

	var version int

	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errMsg.ErrNotFound
		}
		return 0, err
	}

	return version, nil
}
