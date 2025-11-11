package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type userCategoryRepo struct {
	db repo.DBExecutor
}

func NewUserCategory(db repo.DBExecutor) UserCategory {
	return &userCategoryRepo{db: db}
}
