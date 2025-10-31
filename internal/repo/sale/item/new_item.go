package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type saleItem struct {
	db repo.DBExecutor
}

// constrói usando *pgxpool.Pool
func NewItemSaleFromPool(pool *pgxpool.Pool) SaleItemRepo {
	return &saleItem{db: pool}
}

// constrói com um DBExecutor customizado (Tx, mock, etc.)
func NewItemSale(db repo.DBExecutor) SaleItemRepo {
	return &saleItem{db: db}
}
