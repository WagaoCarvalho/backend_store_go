package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

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

// Connect inicializa a conexão com o banco de dados
func Connect(pool PgxPool) (*pgxpool.Pool, error) {
	// Obtém a configuração do banco de dados
	dbConfig := config.LoadDatabaseConfig()

	// Verifica se a URL de conexão foi carregada corretamente
	if dbConfig.ConnURL == "" {
		return nil, errors.New("variável de ambiente DB_CONN_URL não definida")
	}

	// Parseia a configuração do pool
	pgxConfig, err := pool.ParseConfig(dbConfig.ConnURL)
	if err != nil {
		return nil, err
	}

	// Cria a conexão com o pool
	dbPool, err = pool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Conectado ao banco de dados com sucesso!")
	return dbPool, nil
}
