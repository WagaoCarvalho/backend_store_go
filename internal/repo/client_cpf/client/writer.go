package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *clientCpfRepo) Create(ctx context.Context, clientCpf *models.ClientCpf) (*models.ClientCpf, error) {
	const query = `
		INSERT INTO clients_cpf (name, email, cpf, description, status, version)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		clientCpf.Name,
		clientCpf.Email,
		clientCpf.CPF,
		clientCpf.Description,
		clientCpf.Status,
		clientCpf.Version,
	).Scan(&clientCpf.ID, &clientCpf.CreatedAt, &clientCpf.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrDBInvalidForeignKey
		}

		if ok, constraint := errMsgPg.IsUniqueViolation(err); ok {
			return nil, fmt.Errorf("%w: %s", errMsg.ErrDuplicate, constraint)
		}

		if errMsgPg.IsCheckViolation(err) {
			return nil, errMsg.ErrInvalidData
		}

		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return clientCpf, nil
}

func (r *clientCpfRepo) Update(ctx context.Context, clientCpf *models.ClientCpf) error {
	const querySelect = `
		SELECT version
		FROM clients_cpf
		WHERE id = $1
	`

	var currentVersion int
	err := r.db.QueryRow(ctx, querySelect, clientCpf.ID).Scan(&currentVersion)

	if errors.Is(err, pgx.ErrNoRows) {
		return errMsg.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	if currentVersion != clientCpf.Version {
		return errMsg.ErrVersionConflict
	}

	const queryUpdate = `
		UPDATE clients_cpf
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			status = $4,
			description = $5,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at, version
	`

	err = r.db.QueryRow(ctx, queryUpdate,
		clientCpf.Name,
		clientCpf.Email,
		clientCpf.CPF,
		clientCpf.Status,
		clientCpf.Description,
		clientCpf.ID,
	).Scan(&clientCpf.UpdatedAt, &clientCpf.Version)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return errMsg.ErrDuplicate
			case "23514":
				return errMsg.ErrInvalidData
			}
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *clientCpfRepo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM clients_cpf WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
