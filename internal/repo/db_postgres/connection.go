package repo

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/config"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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

func Connect(pool PgxPool) (*pgxpool.Pool, error) {
	dbConfig := config.LoadDatabaseConfig()

	if dbConfig.ConnURL == "" {
		return nil, errMsg.ErrDBConnURLNotDefined
	}

	pgxConfig, err := pool.ParseConfig(dbConfig.ConnURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrDBParseConfig, err)
	}

	dbPool, err = pool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrDBNewPool, err)
	}

	if _, ok := pool.(*RealPgxPool); ok {
		if err := dbPool.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrDBPing, err)
		}
	}

	return dbPool, nil
}
