package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type saleFilterRepo struct {
	db repo.DBExecutor
}

func NewFilterSale(db repo.DBExecutor) SaleFilter {
	return &saleFilterRepo{db: db}
}
