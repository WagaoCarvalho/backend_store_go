package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	var called bool

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler := Logging(next)
	handler.ServeHTTP(rr, req)

	assert.True(t, called, "O handler interno deveria ter sido chamado")
	assert.Equal(t, http.StatusOK, rr.Code)
}
