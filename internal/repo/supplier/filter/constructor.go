package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type supplierFilterRepo struct {
	db repo.DBExecutor
}

func NewFilterSupplier(db repo.DBExecutor) SupplierFilter {
	return &supplierFilterRepo{db: db}
}
