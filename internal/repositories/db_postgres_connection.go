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

func Connect() *pgxpool.Pool {
	url := configs.Database.ConnURL

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Erro ao configurar a conexão: %v\n", err)
	}

	dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v\n", err)
	}

	fmt.Println("Banco de dados conectado com sucesso!")
	return dbPool
}

func TestConnection() {
	con := Connect()
	defer con.Close()

	err := con.Ping(context.Background())
	if err != nil {
		log.Fatalf("Erro ao conectar: %v\n", err)
	}
	fmt.Println("Database connected!")
}
