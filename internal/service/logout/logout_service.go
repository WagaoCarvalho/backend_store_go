package auth

import (
	"context"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type LogoutService interface {
	Logout(ctx context.Context, token string) error
}

type TokenBlacklist interface {
	Add(ctx context.Context, token string, duration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type JWTService interface {
	Parse(tokenString string) (*jwt.Token, error)
	GetExpiration(token *jwt.Token) (time.Duration, error)
}

type logoutService struct {
	blacklist  TokenBlacklist
	logger     *logger.LoggerAdapter
	jwtService JWTService
}

func NewLogoutService(
	blacklist TokenBlacklist,
	logger *logger.LoggerAdapter,
	jwtService JWTService,
) *logoutService {
	return &logoutService{
		blacklist:  blacklist,
		logger:     logger,
		jwtService: jwtService,
	}
}

func (s *logoutService) Logout(ctx context.Context, tokenString string) error {
	const ref = "[logoutService - Logout] - "

	s.logger.Info(ctx, ref+logger.LogLogoutInit, map[string]any{
		"token": tokenString,
	})

	token, err := s.jwtService.Parse(tokenString)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogTokenValidationFail, nil)
		return ErrTokenValidation
	}

	if !token.Valid {
		s.logger.Warn(ctx, ref+logger.LogTokenInvalid, nil)
		return ErrInvalidToken
	}

	expiration, err := s.jwtService.GetExpiration(token)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogClaimExpInvalid, nil)
		return ErrClaimExpInvalid
	}

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
