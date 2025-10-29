package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sale struct {
	db repo.DBExecutor
}

// NewSaleFromPool constrói usando *pgxpool.Pool
func NewSaleFromPool(pool *pgxpool.Pool) SaleRepo {
	return &sale{db: pool}
}

// NewSale constrói com um DBExecutor customizado (Tx, mock, etc.)
func NewSale(db repo.DBExecutor) SaleRepo {
	return &sale{db: db}
}
