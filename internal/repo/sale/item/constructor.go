package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
)

type saleItemRepo struct {
	db repo.DBExecutor
}

func NewItemSale(db repo.DBExecutor) SaleItemRepo {
	return &saleItemRepo{db: db}
}
