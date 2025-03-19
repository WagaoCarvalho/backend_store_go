package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Database Database
	Server   Server
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Não foi possível carregar o .env. Usando variáveis de ambiente.")
	}

	return Config{
		Database: LoadDatabaseConfig(),
		Server:   LoadServerConfig(),
	}
}
