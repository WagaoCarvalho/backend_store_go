package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	context "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/context_utils"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware_GeraENaoSobrescreve(t *testing.T) {
	t.Run("Deve gerar um novo request_id quando não há header", func(t *testing.T) {
		var capturedRequestID string

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedRequestID = context.GetRequestID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler := RequestIDMiddleware()(next)
		handler.ServeHTTP(rr, req)

		responseHeaderID := rr.Header().Get("X-Request-ID")

		assert.NotEmpty(t, capturedRequestID)
		assert.Equal(t, responseHeaderID, capturedRequestID)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Deve manter o request_id original se presente no header", func(t *testing.T) {
		expectedID := "meu-request-id-123"
		var capturedRequestID string

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedRequestID = context.GetRequestID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", expectedID)
		rr := httptest.NewRecorder()

		handler := RequestIDMiddleware()(next)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, expectedID, capturedRequestID)
		assert.Equal(t, expectedID, rr.Header().Get("X-Request-ID"))
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
