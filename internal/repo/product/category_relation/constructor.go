package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"

type productCategoryRelation struct {
	db repo.DBExecutor
}

func NewProductCategoryRelation(db repo.DBExecutor) ProductCategoryRelation {
	return &productCategoryRelation{db: db}
}
