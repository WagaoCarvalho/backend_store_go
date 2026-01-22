package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *clientCpfRepo) GetByID(ctx context.Context, id int64) (*models.ClientCpf, error) {
	const query = `
		SELECT
			id,
			name,
			email,
			cpf,
			description,
			status,
			version,
			created_at,
			updated_at
		FROM clients_cpf
		WHERE id = $1
		LIMIT 1
	`

	clientCpf := &models.ClientCpf{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&clientCpf.ID,
		&clientCpf.Name,
		&clientCpf.Email,
		&clientCpf.CPF,
		&clientCpf.Description,
		&clientCpf.Status,
		&clientCpf.Version,
		&clientCpf.CreatedAt,
		&clientCpf.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w", err)
	}

	return clientCpf, nil
}
