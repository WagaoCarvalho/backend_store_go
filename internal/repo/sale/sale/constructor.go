package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
)

type sale struct {
	db repo.DBExecutor
}

func NewSale(db repo.DBExecutor) SaleRepo {
	return &sale{db: db}
}
