package repo

import (
	"context"
	"fmt"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *clientRepo) GetAll(ctx context.Context, filter *model.ClientFilter) ([]*model.Client, error) {
	query := `
		SELECT id, name, email, cpf, cnpj, description, status, version, created_at, updated_at
		FROM clients
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Name)
		argPos++
	}
	if filter.Email != "" {
		query += fmt.Sprintf(" AND email ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Email)
		argPos++
	}
	if filter.CPF != "" {
		query += fmt.Sprintf(" AND cpf = $%d", argPos)
		args = append(args, filter.CPF)
		argPos++
	}
	if filter.CNPJ != "" {
		query += fmt.Sprintf(" AND cnpj = $%d", argPos)
		args = append(args, filter.CNPJ)
		argPos++
	}
	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filter.Status)
		argPos++
	}
	if filter.Version != nil {
		query += fmt.Sprintf(" AND version = $%d", argPos)
		args = append(args, *filter.Version)
		argPos++
	}
	if filter.CreatedFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *filter.CreatedFrom)
		argPos++
	}
	if filter.CreatedTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *filter.CreatedTo)
		argPos++
	}
	if filter.UpdatedFrom != nil {
		query += fmt.Sprintf(" AND updated_at >= $%d", argPos)
		args = append(args, *filter.UpdatedFrom)
		argPos++
	}
	if filter.UpdatedTo != nil {
		query += fmt.Sprintf(" AND updated_at <= $%d", argPos)
		args = append(args, *filter.UpdatedTo)
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var clients []*model.Client
	for rows.Next() {
		var c model.Client
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Email, &c.CPF, &c.CNPJ,
			&c.Description, &c.Status, &c.Version,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		clients = append(clients, &c)
	}

	// 🔍 Checagem obrigatória de erro de iteração
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return clients, nil
}
