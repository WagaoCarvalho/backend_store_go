package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"

type clientCredit struct {
	db repo.DBExecutor
}

func NewClientCredit(db repo.DBExecutor) ClientCredit {
	return &clientCredit{db: db}
}
