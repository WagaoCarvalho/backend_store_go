package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleRepository interface {
	Create(ctx context.Context, sale *models.Sale) (*models.Sale, error)
	CreateTx(ctx context.Context, tx pgx.Tx, sale *models.Sale) (*models.Sale, error)
	GetByID(ctx context.Context, id int64) (*models.Sale, error)
	GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	Update(ctx context.Context, sale *models.Sale) error
	Delete(ctx context.Context, id int64) error
}

type saleRepository struct {
	db *pgxpool.Pool
}

func NewSaleRepository(db *pgxpool.Pool) SaleRepository {
	return &saleRepository{db: db}
}

func (r *saleRepository) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	const query = `
		INSERT INTO sales (
			client_id, user_id, sale_date, total_amount, total_discount,
			payment_type, status, notes, version, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1, NOW(), NOW())
		RETURNING id, version, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		sale.ClientID,
		sale.UserID,
		sale.SaleDate,
		sale.TotalAmount,
		sale.TotalDiscount,
		sale.PaymentType,
		sale.Status,
		sale.Notes,
	).Scan(&sale.ID, &sale.Version, &sale.CreatedAt, &sale.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return sale, nil
}

func (r *saleRepository) CreateTx(ctx context.Context, tx pgx.Tx, sale *models.Sale) (*models.Sale, error) {
	const query = `
		INSERT INTO sales (
			client_id, user_id, sale_date, total_amount, total_discount,
			payment_type, status, notes, version, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1, NOW(), NOW())
		RETURNING id, version, created_at, updated_at;
	`

	err := tx.QueryRow(ctx, query,
		sale.ClientID,
		sale.UserID,
		sale.SaleDate,
		sale.TotalAmount,
		sale.TotalDiscount,
		sale.PaymentType,
		sale.Status,
		sale.Notes,
	).Scan(&sale.ID, &sale.Version, &sale.CreatedAt, &sale.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return sale, nil
}

func (r *saleRepository) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
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

func (r *saleRepository) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	return r.listByField(ctx, "client_id", clientID, limit, offset, orderBy, orderDir)
}

func (r *saleRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	return r.listByField(ctx, "user_id", userID, limit, offset, orderBy, orderDir)
}

func (r *saleRepository) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	return r.listByField(ctx, "status", status, limit, offset, orderBy, orderDir)
}

func (r *saleRepository) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
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

	var sales []*models.Sale
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
		sales = append(sales, &sale)
	}

	return sales, nil
}

func (r *saleRepository) Update(ctx context.Context, sale *models.Sale) error {
	const query = `
		UPDATE sales
		SET 
			client_id     = $1,
			user_id       = $2,
			sale_date     = $3,
			total_amount  = $4,
			total_discount= $5,
			payment_type  = $6,
			status        = $7,
			notes         = $8,
			version       = version + 1,
			updated_at    = NOW()
		WHERE id = $9 AND version = $10
		RETURNING version, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		sale.ClientID,
		sale.UserID,
		sale.SaleDate,
		sale.TotalAmount,
		sale.TotalDiscount,
		sale.PaymentType,
		sale.Status,
		sale.Notes,
		sale.ID,
		sale.Version,
	).Scan(&sale.Version, &sale.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		if errMsgPg.IsForeignKeyViolation(err) {
			return errMsg.ErrInvalidForeignKey
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *saleRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM sales 
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

// --- Helpers ---
// Reuso para consultas com WHERE <field> = $1
func (r *saleRepository) listByField(ctx context.Context, field string, value interface{}, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, client_id, user_id, sale_date, total_amount, total_discount,
			payment_type, status, notes, version, created_at, updated_at
		FROM sales
		WHERE %s = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3;
	`, field, sanitizeOrderBy(orderBy), sanitizeOrderDir(orderDir))

	rows, err := r.db.Query(ctx, query, value, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var sales []*models.Sale
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
		sales = append(sales, &sale)
	}

	return sales, nil
}

// Sanitize order by to avoid SQL injection
func sanitizeOrderBy(orderBy string) string {
	switch orderBy {
	case "sale_date", "created_at", "updated_at", "total_amount", "status":
		return orderBy
	default:
		return "sale_date"
	}
}

func sanitizeOrderDir(orderDir string) string {
	if orderDir == "ASC" {
		return "ASC"
	}
	return "DESC"
}
