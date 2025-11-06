package repo

import "github.com/jackc/pgx/v5/pgxpool"

type contact struct {
	db *pgxpool.Pool
}

func NewContact(db *pgxpool.Pool) Contact {
	return &contact{db: db}
}
