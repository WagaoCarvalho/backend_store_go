package repo

import "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"

type supplier struct {
	db repo.DBExecutor
}

func NewSupplier(db repo.DBExecutor) Supplier {
	return &supplier{db: db}
}
