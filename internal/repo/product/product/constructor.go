package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
)

type product struct {
	db repo.DBExecutor
}

func NewProduct(db repo.DBExecutor) Product {
	return &product{db: db}
}
