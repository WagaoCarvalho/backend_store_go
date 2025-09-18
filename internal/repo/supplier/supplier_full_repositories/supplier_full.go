package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierFullRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, supplier *models.Supplier) (*models.Supplier, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type supplierFullRepository struct {
	db *pgxpool.Pool
}

func NewSupplierFullRepository(db *pgxpool.Pool) SupplierFullRepository {
	return &supplierFullRepository{
		db: db,
	}
}

func (r *supplierFullRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *supplierFullRepository) CreateTx(ctx context.Context, tx pgx.Tx, supplier *models.Supplier) (*models.Supplier, error) {
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
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return supplier, nil
}
