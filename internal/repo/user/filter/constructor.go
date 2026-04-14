package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type userFilterRepo struct {
	db repo.DBExecutor
}

func NewUserFilter(db repo.DBExecutor) UserFilter {
	return &userFilterRepo{db: db}
}
