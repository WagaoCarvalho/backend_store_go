package repo

import "github.com/jackc/pgx/v5/pgxpool"

type address struct {
	db *pgxpool.Pool
}

func NewAddress(db *pgxpool.Pool) Address {
	return &address{db: db}
}
