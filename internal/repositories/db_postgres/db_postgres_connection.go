package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool
var configs = config.LoadConfig()

type PgxPool interface {
	ParseConfig(connString string) (*pgxpool.Config, error)
	NewWithConfig(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error)
}

type RealPgxPool struct{}

func (r *RealPgxPool) ParseConfig(connString string) (*pgxpool.Config, error) {
	return pgxpool.ParseConfig(connString)
}

func (r *RealPgxPool) NewWithConfig(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
	return pgxpool.NewWithConfig(ctx, config)
}

func Connect(pgx PgxPool) (*pgxpool.Pool, error) {
	url := configs.Database.ConnURL

	config, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao configurar a conexão: %v", err)
	}

	dbPool, err = pgx.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}

	log.Println("Banco de dados conectado com sucesso!")
	return dbPool, nil
}
