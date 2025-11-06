package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"

type userContactRelation struct {
	db repo.DBExecutor
}

func NewUserContactRelation(db repo.DBExecutor) UserContactRelation {
	return &userContactRelation{db: db}
}
