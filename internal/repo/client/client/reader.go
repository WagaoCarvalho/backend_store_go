package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *client) GetByID(ctx context.Context, id int64) (*models.Client, error) {
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

func (r *client) GetByName(ctx context.Context, name string) ([]*models.Client, error) {
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

func (r *client) GetVersionByID(ctx context.Context, id int64) (int, error) {
	const query = `SELECT version FROM clients WHERE id = $1`
	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return version, nil
}

func (r *client) ClientExists(ctx context.Context, clientID int64) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM clients WHERE id=$1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, clientID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return exists, nil
}
