package repo

import "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"

type userCategoryRelation struct {
	db repo.DBExecutor
}

func NewUserCategoryRelation(db repo.DBExecutor) UserCategoryRelation {
	return &userCategoryRelation{db: db}
}
