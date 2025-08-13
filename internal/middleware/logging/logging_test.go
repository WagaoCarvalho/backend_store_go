package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	infoCalled bool
	lastCtx    context.Context
	lastMsg    string
	lastFields map[string]interface{}
}

func (m *mockLogger) Info(ctx context.Context, msg string, extraFields map[string]interface{}) {
	m.infoCalled = true
	m.lastCtx = ctx
	m.lastMsg = msg
	m.lastFields = extraFields
}

func (m *mockLogger) Error(ctx context.Context, err error, msg string, extraFields map[string]interface{}) {
	// Pode implementar se precisar para testes futuros
}

func (m *mockLogger) Warn(ctx context.Context, msg string, extraFields map[string]interface{}) {
	//m.warnCalled = true
	m.lastCtx = ctx
	m.lastMsg = msg
	m.lastFields = extraFields
}

func TestLoggingMiddleware(t *testing.T) {
	mockLog := &mockLogger{}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler := LoggingMiddleware(mockLog)(next)
	handler.ServeHTTP(rr, req)

	assert.True(t, mockLog.infoCalled, "Logger.Info deveria ser chamado")
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "[*** - Request concluída - ***]", mockLog.lastMsg)
	assert.Equal(t, "/test", mockLog.lastFields["path"])
	assert.Equal(t, http.MethodGet, mockLog.lastFields["method"])
	assert.Equal(t, http.StatusOK, mockLog.lastFields["status"])
}
