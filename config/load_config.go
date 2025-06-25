package config

import (
	"os"
	"strings"
)

func LoadAppConfig() App {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // valor padrão
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return App{
		Env:      strings.ToLower(env),
		LogLevel: strings.ToLower(logLevel),
	}
}
