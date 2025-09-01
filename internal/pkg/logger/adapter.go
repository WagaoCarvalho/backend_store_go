package logger

import (
	"context"

	contextUtils "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/context_utils"
	"github.com/sirupsen/logrus"
)

type LogAdapter struct {
	base *logrus.Logger
}

var _ LogAdapterInterface = (*LogAdapter)(nil)

type LogAdapterInterface interface {
	Warn(ctx context.Context, msg string, extraFields map[string]any)
	Info(ctx context.Context, msg string, extraFields map[string]any)
	Error(ctx context.Context, err error, msg string, extraFields map[string]any)
}

func NewLoggerAdapter(base *logrus.Logger) *LogAdapter {
	if base == nil {
		panic("logrus.Logger não pode ser nil")
	}
	return &LogAdapter{base: base}
}

func (l *LogAdapter) WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{
		"request_id": contextUtils.GetRequestID(ctx),
	}

	if uid := contextUtils.GetUserID(ctx); uid != "" {
		fields["system_user_id"] = uid
	}

	return l.base.WithFields(fields)
}

func (l *LogAdapter) Warn(ctx context.Context, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).Warn(msg)
}

func (l *LogAdapter) Info(ctx context.Context, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).Info(msg)
}

func (l *LogAdapter) Error(ctx context.Context, err error, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).WithError(err).Error(msg)
}
