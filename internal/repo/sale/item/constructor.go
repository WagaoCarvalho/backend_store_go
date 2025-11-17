package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
)

type itemSaleRepo struct {
	db repo.DBExecutor
}

func NewItemSale(db repo.DBExecutor) SaleItemRepo {
	return &itemSaleRepo{db: db}
}
