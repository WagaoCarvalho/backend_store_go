package repo

import (
	"context"
	"fmt"
	"strings"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

var allowedSupplierSortFields = map[string]string{
	"id":         "id",
	"name":       "name",
	"cpf":        "cpf",
	"cnpj":       "cnpj",
	"status":     "status",
	"version":    "version",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (r *supplierFilterRepo) Filter(
	ctx context.Context,
	filter *filter.SupplierFilter,
) ([]*model.Supplier, error) {

	base := filter.BaseFilter.WithDefaults()

	query := `
		SELECT
			id,
			name,
			cpf,
			cnpj,
			description,
			status,
			version,
			created_at,
			updated_at
		FROM suppliers
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+filter.Name+"%")
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
	if v, ok := allowedSupplierSortFields[strings.ToLower(base.SortBy)]; ok {
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

	var suppliers []*model.Supplier
	for rows.Next() {
		var s model.Supplier
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.CPF,
			&s.CNPJ,
			&s.Description,
			&s.Status,
			&s.Version,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		suppliers = append(suppliers, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return suppliers, nil
}
