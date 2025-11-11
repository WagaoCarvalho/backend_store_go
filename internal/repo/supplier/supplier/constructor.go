package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type supplierRepo struct {
	db repo.DBExecutor
}

func NewSupplier(db repo.DBExecutor) Supplier {
	return &supplierRepo{db: db}
}
