package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	t.Run("should add CORS headers for normal request", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler := CORS(next)
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
		assert.Contains(t, res.Header.Get("Access-Control-Allow-Methods"), "GET")
		assert.Contains(t, res.Header.Get("Access-Control-Allow-Headers"), "Authorization")
		assert.True(t, nextCalled)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should respond 200 for OPTIONS and not call next", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()

		handler := CORS(next)
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.False(t, nextCalled)

		assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
		assert.Contains(t, res.Header.Get("Access-Control-Allow-Methods"), "DELETE")
		assert.Contains(t, res.Header.Get("Access-Control-Allow-Headers"), "Content-Type")
	})
}
