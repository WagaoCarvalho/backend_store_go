package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type supplierCategoryRelationRepo struct {
	db repo.DBExecutor
}

func NewSupplierCategoryRelation(db repo.DBExecutor) SupplierCategoryRelation {
	return &supplierCategoryRelationRepo{db: db}
}
