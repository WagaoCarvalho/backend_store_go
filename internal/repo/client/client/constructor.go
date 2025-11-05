package repo

import "github.com/jackc/pgx/v5/pgxpool"

type client struct {
	db *pgxpool.Pool
}

func NewClient(db *pgxpool.Pool) Client {
	return &client{db: db}
}
