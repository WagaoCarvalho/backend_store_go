package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type userContactRelationRepo struct {
	db repo.DBExecutor
}

func NewUserContactRelation(db repo.DBExecutor) UserContactRelation {
	return &userContactRelationRepo{db: db}
}
