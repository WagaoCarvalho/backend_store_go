package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

type LogoutService interface {
	Logout(ctx context.Context, token string) error
}

type TokenBlacklist interface {
	Add(ctx context.Context, token string, duration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type logoutService struct {
	blacklist TokenBlacklist
	logger    *logger.LoggerAdapter
	secretKey string
}

func NewLogoutService(blacklist TokenBlacklist, logger *logger.LoggerAdapter, secretKey string) *logoutService {
	return &logoutService{
		blacklist: blacklist,
		logger:    logger,
		secretKey: secretKey,
	}
}

func (s *logoutService) Logout(ctx context.Context, tokenString string) error {
	const ref = "[logoutService - Logout] - "

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Logamos diretamente aqui, pois é erro de segurança
			s.logger.Error(ctx, nil, ref+"método de assinatura inválido", map[string]any{
				"alg": token.Header["alg"],
			})
			return nil, fmt.Errorf("método de assinatura inválido")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		s.logger.Error(ctx, err, ref+"falha ao validar token", nil)
		return fmt.Errorf("erro ao validar token: %w", err)
	}
	if !token.Valid {
		s.logger.Warn(ctx, ref+"token inválido", nil)
		return fmt.Errorf("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error(ctx, nil, ref+"não foi possível converter claims", nil)
		return fmt.Errorf("não foi possível obter claims")
	}

	expUnix, ok := claims["exp"].(float64)
	if !ok {
		s.logger.Error(ctx, nil, ref+"claim 'exp' ausente ou inválida", map[string]any{
			"claims": claims,
		})
		return fmt.Errorf("claim 'exp' ausente ou inválida")
	}

	expiration := time.Until(time.Unix(int64(expUnix), 0))
	if expiration <= 0 {
		s.logger.Warn(ctx, ref+"token já expirado", nil)
		return fmt.Errorf("token já expirado")
	}

	err = s.blacklist.Add(ctx, tokenString, expiration)
	if err != nil {
		s.logger.Error(ctx, err, ref+"erro ao adicionar token à blacklist", map[string]any{
			"token": tokenString,
		})
		return fmt.Errorf("erro ao realizar logout: %w", err)
	}

	s.logger.Info(ctx, ref+"logout realizado com sucesso", map[string]any{
		"token": tokenString,
	})
	return nil
}
