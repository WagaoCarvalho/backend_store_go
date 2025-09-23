package logger

import (
	"context"

	contextUtils "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/context_utils"
	"github.com/sirupsen/logrus"
)

type LogAdapter struct {
	base      *logrus.Logger
	component string
}

// NewLoggerAdapter cria um adapter de log. Por padrão o component é "app".
// Se quiser mudar (ex: "system"), basta passar como segundo argumento.
func NewLoggerAdapter(base *logrus.Logger, component ...string) *LogAdapter {
	if base == nil {
		panic("logrus.Logger não pode ser nil")
	}

	comp := "app"
	if len(component) > 0 && component[0] != "" {
		comp = component[0]
	}

	return &LogAdapter{base: base, component: comp}
}

func (l *LogAdapter) WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{
		"component": l.component,
	}

	// Só adiciona request_id/system_user_id se for app
	if l.component == "app" {
		if rid := contextUtils.GetRequestID(ctx); rid != "" {
			fields["request_id"] = rid
		}
		if uid := contextUtils.GetUserID(ctx); uid != "" {
			fields["system_user_id"] = uid
		}
	}

	return l.base.WithFields(fields)
}

func (l *LogAdapter) Info(ctx context.Context, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).Info(msg)
}

func (l *LogAdapter) Warn(ctx context.Context, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).Warn(msg)
}

func (l *LogAdapter) Error(ctx context.Context, err error, msg string, extraFields map[string]any) {
	l.WithContext(ctx).WithFields(logrus.Fields(extraFields)).WithError(err).Error(msg)
}
