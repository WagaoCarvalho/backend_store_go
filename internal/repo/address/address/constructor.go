package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type addressRepo struct {
	db repo.DBExecutor
}

func NewAddress(db repo.DBExecutor) Address {
	return &addressRepo{db: db}
}
