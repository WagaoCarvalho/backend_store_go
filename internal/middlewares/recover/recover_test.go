package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

// Mock simples do LoggerAdapter só para capturar chamadas do método Error
type mockLoggerAdapter struct {
	errorCalled bool
	lastCtx     context.Context
	lastErr     error
	lastMsg     string
	lastFields  map[string]interface{}
}

func (m *mockLoggerAdapter) WithContext(ctx context.Context) *logrus.Entry {
	return nil // não usado no teste
}

func (m *mockLoggerAdapter) Info(ctx context.Context, msg string, extraFields map[string]interface{}) {
	// Não usado aqui
}

func (m *mockLoggerAdapter) Warn(ctx context.Context, msg string, extraFields map[string]interface{}) {
	//m.warnCalled = true
	m.lastCtx = ctx
	m.lastMsg = msg
	m.lastFields = extraFields
}

func (m *mockLoggerAdapter) Error(ctx context.Context, err error, msg string, extraFields map[string]interface{}) {
	m.errorCalled = true
	m.lastCtx = ctx
	m.lastErr = err
	m.lastMsg = msg
	m.lastFields = extraFields
}

func TestRecoverMiddleware(t *testing.T) {
	mockLog := &mockLoggerAdapter{}

	t.Run("Deve retornar 500 quando ocorrer panic", func(t *testing.T) {
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("algo deu errado")
		})

		handler := RecoverMiddleware(mockLog)(panicHandler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Erro interno do servidor")

		assert.True(t, mockLog.errorCalled, "Logger.Error deveria ser chamado")
		assert.NotNil(t, mockLog.lastErr)
		assert.Equal(t, "Panic capturado", mockLog.lastMsg)
	})

	t.Run("Deve continuar normalmente quando não houver panic", func(t *testing.T) {
		mockLog.errorCalled = false // reset
		called := false
		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Tudo certo"))
		})

		handler := RecoverMiddleware(mockLog)(normalHandler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Tudo certo", rr.Body.String())
		assert.False(t, mockLog.errorCalled, "Logger.Error não deve ser chamado")
	})
}
