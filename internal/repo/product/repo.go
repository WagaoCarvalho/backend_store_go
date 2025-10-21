package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db DBExecutor
}

func NewProductRepositoryFromPool(pool *pgxpool.Pool) ProductRepo {
	return &productRepository{db: pool}
}

func NewProductRepository(db DBExecutor) ProductRepo {
	return &productRepository{db: db}
}
