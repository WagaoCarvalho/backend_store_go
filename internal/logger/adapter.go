package logger

import (
	"context"

	contextutils "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/request/context_utils"
	"github.com/sirupsen/logrus"
)

type LoggerAdapter struct {
	base *logrus.Logger
}

func NewLoggerAdapter(base *logrus.Logger) *LoggerAdapter {
	return &LoggerAdapter{base: base}
}

func (l *LoggerAdapter) WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{
		"request_id": contextutils.GetRequestID(ctx),
	}

	if uid := contextutils.GetUserID(ctx); uid != "" {
		fields["user_id"] = uid
	}

	return l.base.WithFields(fields)
}

func (l *LoggerAdapter) Warn(ctx context.Context, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).Warn(msg)
}

func (l *LoggerAdapter) Info(ctx context.Context, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).Info(msg)
}

func (l *LoggerAdapter) Error(ctx context.Context, err error, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).WithError(err).Error(msg)
}
