package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Server Server
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	return Config{
		Server: LoadServerConfig(),
	}
}
