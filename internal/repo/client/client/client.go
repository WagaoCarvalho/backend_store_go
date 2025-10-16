package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) (*models.Client, error)
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	GetByName(ctx context.Context, name string) ([]*models.Client, error)
	GetVersionByID(ctx context.Context, id int64) (int, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
	ClientExists(ctx context.Context, clientID int64) (bool, error)
}

type clientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
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

func (r *clientRepository) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	const query = `
		SELECT id, name, email, cpf, cnpj, description, status, created_at, updated_at
		FROM clients
		WHERE id = $1
	`
	client := &models.Client{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.CPF,
		&client.CNPJ,
		&client.Description,
		&client.Status,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return client, nil
}

func (r *clientRepository) GetByName(ctx context.Context, name string) ([]*models.Client, error) {
	const query = `
		SELECT id, name, email, cpf, cnpj, description, status, created_at, updated_at
		FROM clients
		WHERE name ILIKE '%' || $1 || '%'
	`
	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var clients []*models.Client
	for rows.Next() {
		c := &models.Client{}
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.CPF,
			&c.CNPJ,
			&c.Description,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func (r *clientRepository) GetVersionByID(ctx context.Context, id int64) (int, error) {
	const query = `SELECT version FROM clients WHERE id = $1`
	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return version, nil
}

func (r *clientRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Client, error) {
	const query = `
		SELECT id, name, email, cpf, cnpj, description, status, created_at, updated_at
		FROM clients
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var clients []*models.Client
	for rows.Next() {
		c := &models.Client{}
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.CPF,
			&c.CNPJ,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return clients, nil
}

func (r *clientRepository) Update(ctx context.Context, client *models.Client) error {
	const query = `
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
		WHERE id = $7 AND version = $8
		RETURNING version
	`

	err := r.db.QueryRow(ctx, query,
		client.Name,
		client.Email,
		client.CPF,
		client.CNPJ,
		client.Status,
		client.Description,
		client.ID,
		client.Version,
	).Scan(&client.Version)

	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.Is(err, pgx.ErrNoRows), errors.Is(err, sql.ErrNoRows):
			// Conflito de versão
			return errMsg.ErrVersionConflict

		case errors.As(err, &pgErr):
			switch pgErr.Code {
			case "23505":
				return errMsg.ErrDuplicate
			case "23514":
				return errMsg.ErrInvalidData
			default:
				return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
			}

		default:
			return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM clients WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}
	return nil
}

func (r *clientRepository) Disable(ctx context.Context, id int64) error {
	const query = `UPDATE clients SET status=false, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}

func (r *clientRepository) Enable(ctx context.Context, id int64) error {
	const query = `UPDATE clients SET status=true, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}

func (r *clientRepository) ClientExists(ctx context.Context, clientID int64) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM clients WHERE id=$1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, clientID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return exists, nil
}
