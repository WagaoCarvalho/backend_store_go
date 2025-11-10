package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
)

type saleItem struct {
	db repo.DBExecutor
}

func NewItemSale(db repo.DBExecutor) SaleItemRepo {
	return &saleItem{db: db}
}
