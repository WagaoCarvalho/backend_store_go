package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

type TokenBlacklist interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

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
				http.Error(w, "Token ausente", http.StatusUnauthorized)
				logger.Warn(r.Context(), ref+"token ausente", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
				logger.Warn(r.Context(), ref+"formato de token inválido", map[string]any{
					"auth_header": authHeader,
				})
				return
			}

			tokenString := parts[1]

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

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					// Logamos erro de segurança e retornamos erro genérico para o handler
					logger.Error(r.Context(), nil, ref+"método de assinatura inválido", map[string]any{
						"alg": token.Header["alg"],
					})
					return nil, fmt.Errorf("Token inválido")
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					http.Error(w, "Token expirado", http.StatusUnauthorized)
					logger.Warn(r.Context(), ref+"token expirado", map[string]any{"token": tokenString})
					return
				case errors.Is(err, jwt.ErrSignatureInvalid):
					http.Error(w, "Assinatura inválida", http.StatusUnauthorized)
					logger.Warn(r.Context(), ref+"assinatura inválida", map[string]any{"token": tokenString})
					return
				default:
					http.Error(w, "Token inválido", http.StatusUnauthorized)
					logger.Warn(r.Context(), ref+"token inválido", map[string]any{"token": tokenString})
					return
				}
			}

			if !token.Valid {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				logger.Warn(r.Context(), ref+"token inválido", map[string]any{
					"token": tokenString,
				})
				return
			}

			ctx := context.WithValue(r.Context(), "user", token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
