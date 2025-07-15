package auth

import (
	"context"
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

	s.logger.Info(ctx, ref+logger.LogLogoutInit, map[string]any{
		"token": tokenString,
	})

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Error(ctx, nil, ref+logger.LogInvalidSigningMethod, map[string]any{
				"alg": token.Header["alg"],
			})
			return nil, ErrInvalidSigningMethod
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogTokenValidationFail, nil)
		return ErrTokenValidation
	}

	if !token.Valid {
		s.logger.Warn(ctx, ref+logger.LogTokenInvalid, nil)
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error(ctx, nil, ref+logger.LogClaimsConversionFail, nil)
		return ErrClaimConversion
	}

	expUnix, ok := claims["exp"].(float64)
	if !ok {
		s.logger.Error(ctx, nil, ref+logger.LogClaimExpInvalid, map[string]any{
			"claims": claims,
		})
		return ErrClaimExpInvalid
	}

	expiration := time.Until(time.Unix(int64(expUnix), 0))
	if expiration <= 0 {
		s.logger.Warn(ctx, ref+logger.LogTokenAlreadyExpired, nil)
		return ErrTokenExpired
	}

	if err := s.blacklist.Add(ctx, tokenString, expiration); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogBlacklistAddFail, map[string]any{
			"token": tokenString,
		})
		return ErrBlacklistAdd
	}

	s.logger.Info(ctx, ref+logger.LogLogoutSuccess, map[string]any{
		"token": tokenString,
	})

	return nil
}
