package repo

import "github.com/jackc/pgx/v5/pgxpool"

type clientCredit struct {
	db *pgxpool.Pool
}

func NewClientCredit(db *pgxpool.Pool) ClientCredit {
	return &clientCredit{db: db}
}
