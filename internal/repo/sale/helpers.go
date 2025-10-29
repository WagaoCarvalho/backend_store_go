package repo

import (
	"context"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *sale) listByField(ctx context.Context, field string, value interface{}, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, client_id, user_id, sale_date, total_amount, total_discount,
			payment_type, status, notes, version, created_at, updated_at
		FROM sales
		WHERE %s = $1
		ORDER BY %s %s
		LIMIT %d OFFSET %d;
	`, field, sanitizeOrderBy(orderBy), sanitizeOrderDir(orderDir), limit, offset)

	rows, err := r.db.Query(ctx, query, value)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	return scanSales(rows)
}

func scanSales(rows pgx.Rows) ([]*models.Sale, error) {
	var result []*models.Sale
	for rows.Next() {
		var sale models.Sale
		if err := rows.Scan(
			&sale.ID,
			&sale.ClientID,
			&sale.UserID,
			&sale.SaleDate,
			&sale.TotalAmount,
			&sale.TotalDiscount,
			&sale.PaymentType,
			&sale.Status,
			&sale.Notes,
			&sale.Version,
			&sale.CreatedAt,
			&sale.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		result = append(result, &sale)
	}
	return result, nil
}

func sanitizeOrderBy(orderBy string) string {
	switch orderBy {
	case "sale_date":
		return "sale_date"
	case "total_amount":
		return "total_amount"
	default:
		return "sale_date"
	}
}

func sanitizeOrderDir(orderDir string) string {
	if strings.ToLower(orderDir) == "desc" {
		return "DESC"
	}
	return "ASC"
}
