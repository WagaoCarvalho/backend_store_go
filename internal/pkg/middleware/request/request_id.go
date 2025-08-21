package middleware

import (
	"net/http"

	contextutils "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/context_utils"
	"github.com/google/uuid"
)

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			ctx := contextutils.SetRequestID(r.Context(), requestID)
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
