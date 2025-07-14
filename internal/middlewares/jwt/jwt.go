package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

type TokenBlacklist interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type contextKey string

const userClaimsKey = contextKey("user")

func IsAuthByBearerToken(
	blacklist TokenBlacklist,
	logger *logger.LoggerAdapter,
	secretKey string,
) func(http.Handler) http.Handler {

	const ref = "[IsAuthByBearerToken] - "

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn(r.Context(), ref+"token ausente", nil)
				http.Error(w, "Token ausente", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn(r.Context(), ref+"formato de token inválido", map[string]any{
					"auth_header": authHeader,
				})
				http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Blacklist
			isRevoked, err := blacklist.IsBlacklisted(r.Context(), tokenString)
			if err != nil {
				logger.Error(r.Context(), err, ref+"erro ao consultar blacklist", map[string]any{
					"token": tokenString,
				})
				http.Error(w, "Erro interno de autenticação", http.StatusInternalServerError)
				return
			}
			if isRevoked {
				logger.Warn(r.Context(), ref+"token revogado", map[string]any{
					"token": tokenString,
				})
				http.Error(w, "Token revogado", http.StatusUnauthorized)
				return
			}

			// Parse do token
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					logger.Error(r.Context(), nil, ref+"método de assinatura inválido", map[string]any{
						"alg": token.Header["alg"],
					})
					return nil, fmt.Errorf("token inválido")
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					logger.Warn(r.Context(), ref+"token expirado", map[string]any{"token": tokenString})
					http.Error(w, "Token expirado", http.StatusUnauthorized)
				case errors.Is(err, jwt.ErrSignatureInvalid):
					logger.Warn(r.Context(), ref+"assinatura inválida", map[string]any{"token": tokenString})
					http.Error(w, "Assinatura inválida", http.StatusUnauthorized)
				default:
					logger.Warn(r.Context(), ref+"token inválido", map[string]any{"token": tokenString})
					http.Error(w, "Token inválido", http.StatusUnauthorized)
				}
				return
			}

			if !token.Valid {
				logger.Warn(r.Context(), ref+"token inválido (parse ok, mas invalid)", map[string]any{"token": tokenString})
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			// Verifica expiração manualmente
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					logger.Warn(r.Context(), ref+"token expirado (manual check)", map[string]any{"token": tokenString})
					http.Error(w, "Token expirado", http.StatusUnauthorized)
					return
				}
			} else {
				logger.Warn(r.Context(), ref+"campo exp ausente ou inválido", nil)
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			// Injeta dados relevantes no contexto
			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
