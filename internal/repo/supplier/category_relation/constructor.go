package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"

type supplierCategoryRelation struct {
	db repo.DBExecutor
}

func NewSupplierCategoryRelation(db repo.DBExecutor) SupplierCategoryRelation {
	return &supplierCategoryRelation{db: db}
}
