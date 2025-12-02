package repo

import (
	"context"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *saleRepo) listByField(
	ctx context.Context,
	field string,
	value any,
	limit, offset int,
	orderBy, orderDir string,
) ([]*models.Sale, error) {

	query := fmt.Sprintf(`
		SELECT 
			id,
			client_id,
			user_id,
			sale_date,
			total_amount,
			total_discount,
			payment_type,
			status,
			notes,
			version,
			created_at,
			updated_at
		FROM sales
		WHERE %s = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3;
	`, sanitizeField(field), sanitizeOrderBy(orderBy), sanitizeOrderDir(orderDir))

	rows, err := r.db.Query(ctx, query, value, limit, offset)
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
			&sale.ClientID, // *int64
			&sale.UserID,   // *int64
			&sale.SaleDate,
			&sale.TotalAmount,
			&sale.TotalSaleDiscount,
			&sale.PaymentType,
			&sale.Status,
			&sale.Notes, // string, NULL vira ""
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

func sanitizeField(field string) string {
	switch field {
	case "client_id", "user_id", "status", "payment_type":
		return field
	default:
		return "client_id"
	}
}

func sanitizeOrderBy(orderBy string) string {
	switch orderBy {
	case "sale_date", "total_amount", "created_at":
		return orderBy
	default:
		return "sale_date"
	}
}

func sanitizeOrderDir(orderDir string) string {
	if strings.ToLower(strings.TrimSpace(orderDir)) == "desc" {
		return "DESC"
	}
	return "ASC"
}
