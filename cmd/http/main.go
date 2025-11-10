package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
	routes "github.com/WagaoCarvalho/backend_store_go/internal/route"
	"github.com/sirupsen/logrus"
)

func main() {
	configs := config.LoadConfig()
	port := configs.Server.Port
	if port == "" {
		port = "5000"
	}

	level := logrus.InfoLevel
	switch configs.App.LogLevel {
	case "debug":
		level = logrus.DebugLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	}

	// Logger para sistema (infra)
	systemRawLogger := logger.NewLogger(logger.LogConfig{
		Environment: configs.App.Env,
		LogFile:     "logs/system.log",
		Level:       level,
	})
	systemLogger := logger.NewLoggerAdapter(systemRawLogger, "system")

	// Logger para app (requests, handlers, etc.)
	appRawLogger := logger.NewLogger(logger.LogConfig{
		Environment: configs.App.Env,
		LogFile:     "logs/app.log",
		Level:       level,
	})
	appLogger := logger.NewLoggerAdapter(appRawLogger)

	// Conecta DB
	db, err := repo.Connect(&repo.RealPgxPool{})
	if err != nil {
		systemLogger.Error(context.TODO(), err, "❌ Erro ao conectar ao banco de dados", nil)
		os.Exit(1)
	}
	defer db.Close()
	systemLogger.Info(context.TODO(), "[✅ - DB CONECTADO -]", nil)

	// Início servidor
	systemLogger.Info(context.TODO(), "[✅ - SERVIDOR INICIADO -]", map[string]any{
		"env":  configs.App.Env,
		"port": port,
	})

	r := routes.NewRouter(appLogger)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// Shutdown gracioso
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-quit
		systemLogger.Info(context.TODO(), "[🔹 - SHUTDOWN INICIADO -]", map[string]any{"signal": sig.String()})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			systemLogger.Error(context.TODO(), err, "❌ Erro durante shutdown", nil)
		} else {
			systemLogger.Info(context.TODO(), "[✅ - SERVIDOR ENCERRADO -]", nil)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		systemLogger.Error(context.TODO(), err, "❌ Falha no ListenAndServe", nil)
	}
}
