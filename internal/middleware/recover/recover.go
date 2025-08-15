package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
)

func RecoverMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(r.Context(), fmt.Errorf("%v", err), "Panic capturado", map[string]any{
						"error":  err,
						"stack":  string(debug.Stack()),
						"path":   r.URL.Path,
						"method": r.Method,
					})

					http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
