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
	Env      string
	LogLevel string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Aviso: Nenhum arquivo .env encontrado. Usando variáveis de ambiente do sistema.")
	}

	return Config{
		Database: LoadDatabaseConfig(),
		Jwt:      LoadJwtConfig(),
		Server:   LoadServerConfig(),
		App:      LoadAppConfig(),
	}
}
