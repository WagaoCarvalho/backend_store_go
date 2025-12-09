package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type clientFilterRepo struct {
	db repo.DBExecutor
}

func NewFilterClient(db repo.DBExecutor) ClientFilter {
	return &clientFilterRepo{db: db}
}
