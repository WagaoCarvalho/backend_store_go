package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/routes"
	"github.com/sirupsen/logrus"
)

func main() {
	// Carrega configurações
	configs := config.LoadConfig()
	port := configs.Server.Port
	if port == "" {
		port = "5000"
	}

	// Define nível de log baseado na configuração
	level := logrus.InfoLevel
	switch configs.App.LogLevel {
	case "debug":
		level = logrus.DebugLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	}

	// Logger base com rotação
	rawLogger := logger.NewLogger(logger.LoggerConfig{
		Environment: configs.App.Env,
		LogFile:     "logs/app.log",
		Level:       level,
	})

	// Logger com suporte a contexto (request_id)
	appLogger := logger.NewLoggerAdapter(rawLogger)

	// Log inicial (sem request_id)
	appLogger.Info(context.TODO(), "[*** - Servidor iniciado - ***]", map[string]any{
		"env":  configs.App.Env,
		"port": port,
	})

	// Inicializa roteador com logger adaptado
	r := routes.NewRouter(appLogger)

	// Inicia o servidor HTTP
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
