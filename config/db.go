package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	ConnURL string
}

func LoadDatabaseConfig() Database {
	// Carrega .env apenas se não estiver rodando no GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Aviso: Nenhum arquivo .env encontrado. Usando variáveis de ambiente do sistema.")
		}
	}

	return Database{
		ConnURL: os.Getenv("DB_CONN_URL"),
	}
}
