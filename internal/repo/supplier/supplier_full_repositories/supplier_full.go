package repositories

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierFullRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, supplier *models.Supplier) (*models.Supplier, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type supplierFullRepository struct {
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewSupplierFullRepository(db *pgxpool.Pool, logger logger.LoggerAdapterInterface) SupplierFullRepository {
	return &supplierFullRepository{
		db:     db,
		logger: logger,
	}
}

func (r *supplierFullRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *supplierFullRepository) CreateTx(ctx context.Context, tx pgx.Tx, supplier *models.Supplier) (*models.Supplier, error) {
	const ref = "[supplierRepository - CreateTx] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name":   supplier.Name,
		"cnpj":   supplier.CNPJ,
		"cpf":    supplier.CPF,
		"status": supplier.Status,
	})

	const query = `
		INSERT INTO suppliers (name, cnpj, cpf, status, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 1, NOW(), NOW())
		RETURNING id, version, created_at, updated_at
	`

	err := tx.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.Status,
	).Scan(
		&supplier.ID,
		&supplier.Version,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name":   supplier.Name,
			"cnpj":   supplier.CNPJ,
			"cpf":    supplier.CPF,
			"status": supplier.Status,
		})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": supplier.ID,
		"name":        supplier.Name,
	})

	return supplier, nil
}
