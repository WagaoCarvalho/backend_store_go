package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *clientRepo) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	const query = `
		INSERT INTO clients (name, email, cpf, cnpj, description, status, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		client.Name,
		client.Email,
		client.CPF,
		client.CNPJ,
		client.Description,
		client.Status,
		client.Version,
	).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return nil, errMsg.ErrDuplicate
			case "23514":
				return nil, errMsg.ErrInvalidData
			}
		}

		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return client, nil
}

func (r *clientRepo) Update(ctx context.Context, client *models.Client) error {

	const querySelect = `
		SELECT version
		FROM clients
		WHERE id = $1
	`

	var currentVersion int
	err := r.db.QueryRow(ctx, querySelect, client.ID).Scan(&currentVersion)

	if errors.Is(err, pgx.ErrNoRows) {
		return errMsg.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: erro ao consultar cliente: %v", errMsg.ErrUpdate, err)
	}

	if currentVersion != client.Version {
		return errMsg.ErrVersionConflict
	}

	const queryUpdate = `
		UPDATE clients
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			cnpj = $4,
			status = $5,
			description = $6,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at, version
	`

	err = r.db.QueryRow(ctx, queryUpdate,
		client.Name,
		client.Email,
		client.CPF,
		client.CNPJ,
		client.Status,
		client.Description,
		client.ID,
	).Scan(&client.UpdatedAt, &client.Version)

	if err != nil {
		var pgErr *pgconn.PgError

		// Erros de banco que fazem sentido manter
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return errMsg.ErrDuplicate
			case "23514":
				return errMsg.ErrInvalidData
			}
		}

		return fmt.Errorf("%w: erro ao atualizar cliente: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *clientRepo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM clients WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
