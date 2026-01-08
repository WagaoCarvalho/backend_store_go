package repo

import (
	"context"
	"fmt"
	"strings"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

var allowedSaleSortFields = map[string]string{
	"id":           "id",
	"sale_date":    "sale_date",
	"total_amount": "total_amount",
	"payment_type": "payment_type",
	"status":       "status",
	"version":      "version",
	"created_at":   "created_at",
	"updated_at":   "updated_at",
}

func (r *saleFilterRepo) Filter(ctx context.Context, filter *filter.SaleFilter) ([]*model.Sale, error) {

	base := filter.BaseFilter.WithDefaults()

	query := `
		SELECT
			id,
			client_id,
			user_id,
			sale_date,
			total_items_amount,
			total_items_discount,
			total_sale_discount,
			total_amount,
			payment_type,
			status,
			notes,
			version,
			created_at,
			updated_at
		FROM sales
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filter.ClientID != nil {
		query += fmt.Sprintf(" AND client_id = $%d", argPos)
		args = append(args, *filter.ClientID)
		argPos++
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *filter.UserID)
		argPos++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filter.Status)
		argPos++
	}

	if filter.PaymentType != "" {
		query += fmt.Sprintf(" AND payment_type = $%d", argPos)
		args = append(args, filter.PaymentType)
		argPos++
	}

	if filter.SaleDateFrom != nil {
		query += fmt.Sprintf(" AND sale_date >= $%d", argPos)
		args = append(args, *filter.SaleDateFrom)
		argPos++
	}

	if filter.SaleDateTo != nil {
		query += fmt.Sprintf(" AND sale_date <= $%d", argPos)
		args = append(args, *filter.SaleDateTo)
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

	// ORDER BY seguro
	sortField := "sale_date"
	if v, ok := allowedSaleSortFields[strings.ToLower(base.SortBy)]; ok {
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

	var sales []*model.Sale
	for rows.Next() {
		var s model.Sale
		if err := rows.Scan(
			&s.ID,
			&s.ClientID,
			&s.UserID,
			&s.SaleDate,
			&s.TotalItemsAmount,
			&s.TotalItemsDiscount,
			&s.TotalSaleDiscount,
			&s.TotalAmount,
			&s.PaymentType,
			&s.Status,
			&s.Notes,
			&s.Version,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		sales = append(sales, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return sales, nil
}
