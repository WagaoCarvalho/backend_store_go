package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type supplierCategoryRepo struct {
	db repo.DBExecutor
}

func NewSupplierCategory(db repo.DBExecutor) SupplierCategory {
	return &supplierCategoryRepo{db: db}
}
