package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type productFilterRepo struct {
	db repo.DBExecutor
}

func NewFilterProduct(db repo.DBExecutor) ProductFilter {
	return &productFilterRepo{db: db}
}
