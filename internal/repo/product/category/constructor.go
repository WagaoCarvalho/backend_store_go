package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type productCategory struct {
	db repo.DBExecutor
}

func NewProductCategory(db repo.DBExecutor) ProductCategory {
	return &productCategory{db: db}
}
