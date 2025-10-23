package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type product struct {
	db repo.DBExecutor
}

func NewProductFromPool(pool *pgxpool.Pool) ProductRepo {
	return &product{db: pool}
}

func NewProduct(db repo.DBExecutor) ProductRepo {
	return &product{db: db}
}
