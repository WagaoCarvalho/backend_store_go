package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type clientCpfRepo struct {
	db repo.DBExecutor
}

func NewClientCpfRepo(db repo.DBExecutor) ClientCpf {
	return &clientCpfRepo{db: db}
}
