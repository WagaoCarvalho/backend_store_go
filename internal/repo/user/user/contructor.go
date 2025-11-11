package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type userRepo struct {
	db repo.DBExecutor
}

func NewUser(db repo.DBExecutor) User {
	return &userRepo{
		db: db,
	}
}
