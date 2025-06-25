package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Database Database
	Jwt      Jwt
	Server   Server
	App      App
}

type App struct {
	Env      string // dev, prod, staging
	LogLevel string // debug, info, warn, error
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Aviso: Nenhum arquivo .env encontrado. Usando variáveis de ambiente do sistema.")
	}

	return Config{
		Database: LoadDatabaseConfig(),
		Jwt:      LoadJwtConfig(),
		Server:   LoadServerConfig(),
		App:      LoadAppConfig(),
	}
}
