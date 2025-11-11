package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type productCategoryRelationRepo struct {
	db repo.DBExecutor
}

func NewProductCategoryRelation(db repo.DBExecutor) ProductCategoryRelationRepo {
	return &productCategoryRelationRepo{db: db}
}
