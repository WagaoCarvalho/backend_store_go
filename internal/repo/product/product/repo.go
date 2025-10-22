package repo

import (
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db repo.DBExecutor
}

func NewProductRepositoryFromPool(pool *pgxpool.Pool) ProductRepo {
	return &productRepository{db: pool}
}

func NewProductRepository(db repo.DBExecutor) ProductRepo {
	return &productRepository{db: db}
}
