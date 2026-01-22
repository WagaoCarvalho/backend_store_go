package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type clientCpfFilterRepo struct {
	db repo.DBExecutor
}

func NewFilterClientCpf(db repo.DBExecutor) ClientCpfFilter {
	return &clientCpfFilterRepo{db: db}
}
