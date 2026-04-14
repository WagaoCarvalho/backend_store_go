package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type addressFilterRepo struct {
	db repo.DBExecutor
}

func NewFilterAddress(db repo.DBExecutor) AddressFilter {
	return &addressFilterRepo{db: db}
}
