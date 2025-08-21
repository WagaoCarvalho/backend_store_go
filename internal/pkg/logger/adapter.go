package logger

import (
	"context"

	contextutils "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/context_utils"
	"github.com/sirupsen/logrus"
)

type LoggerAdapter struct {
	base *logrus.Logger
}

var _ LoggerAdapterInterface = (*LoggerAdapter)(nil)

type LoggerAdapterInterface interface {
	Warn(ctx context.Context, msg string, extraFields map[string]any)
	Info(ctx context.Context, msg string, extraFields map[string]any)
	Error(ctx context.Context, err error, msg string, extraFields map[string]any)
}

func NewLoggerAdapter(base *logrus.Logger) *LoggerAdapter {
	if base == nil {
		panic("logrus.Logger não pode ser nil")
	}
	return &LoggerAdapter{base: base}
}

func (l *LoggerAdapter) WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{
		"request_id": contextutils.GetRequestID(ctx),
	}

	if uid := contextutils.GetUserID(ctx); uid != "" {
		fields["system_user_id"] = uid
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
