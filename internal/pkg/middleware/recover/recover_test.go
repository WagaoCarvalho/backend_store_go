package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

type mockLoggerAdapter struct {
	errorCalled bool
	lastCtx     context.Context
	lastErr     error
	lastMsg     string
	lastFields  map[string]interface{}
}

func (m *mockLoggerAdapter) WithContext(_ context.Context) *logrus.Entry {
	return nil
}

func (m *mockLoggerAdapter) Info(_ context.Context, _ string, _ map[string]interface{}) {

}

func (m *mockLoggerAdapter) Warn(ctx context.Context, msg string, extraFields map[string]interface{}) {

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
		panicHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
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
		mockLog.errorCalled = false
		called := false
		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("Tudo certo")); err != nil {
				// você pode logar ou ignorar, dependendo do teste
				t.Errorf("erro ao escrever resposta: %v", err)
			}
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
