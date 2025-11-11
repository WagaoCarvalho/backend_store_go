package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type contactRepo struct {
	db repo.DBExecutor
}

func NewContact(db repo.DBExecutor) Contact {
	return &contactRepo{db: db}
}
