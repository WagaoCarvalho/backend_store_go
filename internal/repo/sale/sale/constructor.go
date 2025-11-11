package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
)

type saleRepo struct {
	db repo.DBExecutor
}

func NewSale(db repo.DBExecutor) SaleRepo {
	return &saleRepo{db: db}
}
