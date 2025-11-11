package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

func (r *saleRepo) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
	const query = `
		SELECT 
			id, client_id, user_id, sale_date, total_amount, total_discount,
			payment_type, status, notes, version, created_at, updated_at
		FROM sales
		WHERE id = $1;
	`

	var sale models.Sale
	err := r.db.QueryRow(ctx, query, id).Scan(
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
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &sale, nil
}

func (r *saleRepo) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	return r.listByField(ctx, "client_id", clientID, limit, offset, orderBy, orderDir)
}

func (r *saleRepo) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	return r.listByField(ctx, "user_id", userID, limit, offset, orderBy, orderDir)
}

func (r *saleRepo) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	return r.listByField(ctx, "status", status, limit, offset, orderBy, orderDir)
}

func (r *saleRepo) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, client_id, user_id, sale_date, total_amount, total_discount,
			payment_type, status, notes, version, created_at, updated_at
		FROM sales
		WHERE sale_date BETWEEN $1 AND $2
		ORDER BY %s %s
		LIMIT $3 OFFSET $4;
	`, sanitizeOrderBy(orderBy), sanitizeOrderDir(orderDir))

	rows, err := r.db.Query(ctx, query, start, end, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	return scanSales(rows)
}
