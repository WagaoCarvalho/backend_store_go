package logger

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Environment string
	LogFile     string
	Level       logrus.Level
}

type Logger interface {
	Info(ctx context.Context, msg string, extraFields map[string]any)
	Warn(ctx context.Context, msg string, extraFields map[string]any)
	Error(ctx context.Context, err error, msg string, extraFields map[string]any)
}

func NewLogger(cfg LogConfig) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(cfg.Level)

	logWriter := &lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30,   // dias
		Compress:   true, // compactar logs antigos
	}

	if cfg.Environment == "prod" {
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetOutput(logWriter) // só arquivo em produção
	} else {
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		// grava no arquivo e no stdout simultaneamente
		log.SetOutput(io.MultiWriter(os.Stdout, logWriter))
	}

	log.SetReportCaller(true)

	return log
}
