package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	middlewares "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/context_utils"
	"github.com/golang-jwt/jwt/v5"
)

type TokenBlacklist interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type contextKey string

const userClaimsKey = contextKey("user")

func IsAuthByBearerToken(
	blacklist TokenBlacklist,
	logger_adapter *logger.LoggerAdapter,
	secretKey string,
) func(http.Handler) http.Handler {

	const ref = "[IsAuthByBearerToken] - "

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenMissing, nil)
				http.Error(w, ErrTokenMissing.Error(), http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenInvalidFormat, map[string]any{
					"auth_header": authHeader,
				})
				http.Error(w, ErrTokenInvalidFormat.Error(), http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			isRevoked, err := blacklist.IsBlacklisted(r.Context(), tokenString)
			if err != nil {
				logger_adapter.Error(r.Context(), err, ref+logger.LogAuthBlacklistError, map[string]any{
					"token": tokenString,
				})
				http.Error(w, ErrInternalAuth.Error(), http.StatusInternalServerError)
				return
			}

			if isRevoked {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenRevoked, map[string]any{
					"token": tokenString,
				})
				http.Error(w, ErrTokenRevoked.Error(), http.StatusUnauthorized)
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					logger_adapter.Error(r.Context(), nil, ref+logger.LogAuthInvalidSigningMethod, map[string]any{
						"alg": token.Header["alg"],
					})
					return nil, ErrInvalidSigningMethod
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenExpired, map[string]any{"token": tokenString})
					http.Error(w, ErrTokenExpired.Error(), http.StatusUnauthorized)
				case errors.Is(err, jwt.ErrSignatureInvalid):
					logger_adapter.Warn(r.Context(), ref+logger.LogAuthInvalidSignature, map[string]any{"token": tokenString})
					http.Error(w, ErrInvalidSignature.Error(), http.StatusUnauthorized)
				default:
					logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenInvalid, map[string]any{"token": tokenString})
					http.Error(w, ErrTokenInvalid.Error(), http.StatusUnauthorized)
				}
				return
			}

			if !token.Valid {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenInvalidParsed, map[string]any{"token": tokenString})
				http.Error(w, ErrTokenInvalid.Error(), http.StatusUnauthorized)
				return
			}

			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					logger_adapter.Warn(r.Context(), ref+logger.LogAuthTokenExpiredManualCheck, map[string]any{"token": tokenString})
					http.Error(w, ErrTokenExpired.Error(), http.StatusUnauthorized)
					return
				}
			} else {
				logger_adapter.Warn(r.Context(), ref+logger.LogAuthExpClaimInvalid, nil)
				http.Error(w, ErrInvalidExpClaim.Error(), http.StatusUnauthorized)
				return
			}

			// Extrai user_id do token
			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				logger_adapter.Warn(r.Context(), ref+"claim user_id ausente ou inválida", nil)
				http.Error(w, ErrTokenInvalid.Error(), http.StatusUnauthorized)
				return
			}

			// Injeta claims e user_id no contexto
			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			ctx = middlewares.SetUserID(ctx, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
