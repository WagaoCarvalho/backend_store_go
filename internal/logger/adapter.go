package logger

import (
	"context"

	contextutils "github.com/WagaoCarvalho/backend_store_go/internal/context_utils"
	"github.com/sirupsen/logrus"
)

type LoggerAdapter struct {
	base *logrus.Logger
}

func NewLoggerAdapter(base *logrus.Logger) *LoggerAdapter {
	return &LoggerAdapter{base: base}
}

func (l *LoggerAdapter) WithContext(ctx context.Context) *logrus.Entry {
	return l.base.WithField("request_id", contextutils.GetRequestID(ctx))
}

func (l *LoggerAdapter) Info(ctx context.Context, msg string, extraFields map[string]interface{}) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).Info(msg)
}

func (l *LoggerAdapter) Error(ctx context.Context, err error, msg string, extraFields map[string]interface{}) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).WithError(err).Error(msg)
}
