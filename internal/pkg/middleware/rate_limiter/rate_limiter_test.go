package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	calledCount := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calledCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimiter(next)

	ip := "127.0.0.1:1234"

	// Até 3 requisições devem ser aceitas
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ip

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Requisição %d deveria passar", i+1)
	}

	// Requisições após o limite devem ser bloqueadas
	for i := 3; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ip

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Requisição %d deveria ser bloqueada", i+1)
	}

	assert.Equal(t, 3, calledCount, "Apenas 3 requisições devem passar pelo middleware")
}
