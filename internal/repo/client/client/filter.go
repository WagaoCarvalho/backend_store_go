package repo

import (
	"context"
	"fmt"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *client) GetAll(ctx context.Context, f *model.ClientFilter) ([]*model.Client, error) {
	query := `
		SELECT id, name, email, cpf, cnpj, description, status, created_at, updated_at
		FROM clients
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if f.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argID)
		args = append(args, "%"+f.Name+"%")
		argID++
	}
	if f.Email != "" {
		query += fmt.Sprintf(" AND email ILIKE $%d", argID)
		args = append(args, "%"+f.Email+"%")
		argID++
	}
	if f.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argID)
		args = append(args, *f.Status)
		argID++
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var clients []*model.Client
	for rows.Next() {
		c := &model.Client{}
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Email, &c.CPF, &c.CNPJ,
			&c.Description, &c.Status, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return clients, nil
}
