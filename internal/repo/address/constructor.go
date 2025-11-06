package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"

type address struct {
	db repo.DBExecutor
}

func NewAddress(db repo.DBExecutor) Address {
	return &address{db: db}
}
