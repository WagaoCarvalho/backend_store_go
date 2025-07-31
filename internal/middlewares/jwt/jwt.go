package middlewares

import (
	"context"
	"net/http"
	"strings"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	middlewares "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/context_utils"
)

type TokenBlacklist interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type contextKey string

const userClaimsKey = contextKey("user")

func IsAuthByBearerToken(
	blacklist TokenBlacklist,
	logger_adapter *logger.LoggerAdapter,
	jwtService auth.JWTService,
) func(http.Handler) http.Handler {
	const ref = "[IsAuthByBearerToken] - "

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenMissing, nil)
				http.Error(w, auth.ErrTokenMissing.Error(), http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenInvalidFormat, map[string]any{
					"auth_header": authHeader,
				})
				http.Error(w, auth.ErrTokenInvalidFormat.Error(), http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			isRevoked, err := blacklist.IsBlacklisted(r.Context(), tokenString)
			if err != nil {
				logger_adapter.Error(r.Context(), err, ref+logger.LogAuthBlacklistError, map[string]any{
					"token": tokenString,
				})
				http.Error(w, auth.ErrInternalAuth.Error(), http.StatusInternalServerError)
				return
			}

			if isRevoked {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenRevoked, map[string]any{
					"token": tokenString,
				})
				http.Error(w, auth.ErrTokenRevoked.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				logger_adapter.Warn(r.Context(), ref+"token inválido", map[string]any{"err": err.Error()})
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				logger_adapter.Warn(r.Context(), ref+"claim user_id ausente ou inválida", nil)
				http.Error(w, auth.ErrTokenInvalid.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			ctx = middlewares.SetUserID(ctx, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
