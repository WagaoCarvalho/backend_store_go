package config

import "os"

type Jwt struct {
	SecretKey string
}

func LoadJwtConfig() Jwt {
	return Jwt{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}
