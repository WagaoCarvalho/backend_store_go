package repo

import (
	"context"
	"fmt"
	"strings"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

var allowedUserSortFields = map[string]string{
	"id":         "id",
	"username":   "username",
	"email":      "email",
	"status":     "status",
	"version":    "version",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (r *userFilterRepo) Filter(
	ctx context.Context,
	filter *filter.UserFilter,
) ([]*model.User, error) {

	base := filter.BaseFilter.WithDefaults()

	query := `
		SELECT
			id,
			username,
			email,
			password_hash,
			description,
			status,
			version,
			created_at,
			updated_at
		FROM users
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filter.Username != "" {
		query += fmt.Sprintf(" AND username ILIKE $%d", argPos)
		args = append(args, "%"+filter.Username+"%")
		argPos++
	}

	if filter.Email != "" {
		query += fmt.Sprintf(" AND email ILIKE $%d", argPos)
		args = append(args, "%"+filter.Email+"%")
		argPos++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filter.Status)
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

	// ORDER BY seguro
	sortField := "created_at"
	if v, ok := allowedUserSortFields[strings.ToLower(base.SortBy)]; ok {
		sortField = v
	}

	sortOrder := strings.ToLower(base.SortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query += fmt.Sprintf(
		" ORDER BY %s %s LIMIT %d OFFSET %d",
		sortField,
		sortOrder,
		base.Limit,
		base.Offset,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	// Inicializa slice vazia para garantir que nunca retorna nil
	users := make([]*model.User, 0)

	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.UID,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.Description,
			&u.Status,
			&u.Version,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return users, nil
}
