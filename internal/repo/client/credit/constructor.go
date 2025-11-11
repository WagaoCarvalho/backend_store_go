package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type clientCreditRepo struct {
	db repo.DBExecutor
}

func NewClientCredit(db repo.DBExecutor) ClientCredit {
	return &clientCreditRepo{db: db}
}
