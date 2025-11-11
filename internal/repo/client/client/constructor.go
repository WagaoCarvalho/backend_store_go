package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type clientRepo struct {
	db repo.DBExecutor
}

func NewClient(db repo.DBExecutor) Client {
	return &clientRepo{db: db}
}
