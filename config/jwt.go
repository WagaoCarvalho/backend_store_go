package config

import (
	"os"
	"strconv"
	"time"
)

type Jwt struct {
	SecretKey     string
	Issuer        string
	Audience      string
	TokenDuration time.Duration
}

func LoadJwtConfig() Jwt {
	durationStr := os.Getenv("JWT_TOKEN_DURATION")
	if durationStr == "" {
		durationStr = "3600" // padrão: 1 hora em segundos
	}

	seconds, err := strconv.Atoi(durationStr)
	if err != nil {
		seconds = 3600
	}

	return Jwt{
		SecretKey:     os.Getenv("JWT_SECRET_KEY"),
		Issuer:        os.Getenv("JWT_ISSUER"),
		Audience:      os.Getenv("JWT_AUDIENCE"),
		TokenDuration: time.Duration(seconds) * time.Second,
	}
}
