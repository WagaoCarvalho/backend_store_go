package auth

import (
	"context"
	"strings"
	"time"

	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/golang-jwt/jwt/v5"
)

type TokenBlacklist interface {
	Add(ctx context.Context, token string, duration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type LogoutService interface {
	Logout(ctx context.Context, token string) error
}

type JWTService interface {
	Parse(tokenString string) (*jwt.Token, error)
	GetExpiration(token *jwt.Token) (time.Duration, error)
}

type logoutService struct {
	blacklist  TokenBlacklist
	jwtService JWTService
}

func NewLogoutService(
	blacklist TokenBlacklist,
	jwtService JWTService,
) LogoutService {
	return &logoutService{
		blacklist:  blacklist,
		jwtService: jwtService,
	}
}

func (s *logoutService) Logout(ctx context.Context, tokenString string) error {
	if strings.TrimSpace(tokenString) == "" {
		return err_msg.ErrInvalidToken
	}

	token, err := s.jwtService.Parse(tokenString)
	if err != nil {
		return err_msg.ErrTokenValidation
	}

	if !token.Valid {
		return err_msg.ErrInvalidToken
	}

	expiration, err := s.jwtService.GetExpiration(token)
	if err != nil {
		return err_msg.ErrClaimExpInvalid
	}

	if expiration <= 0 {
		return err_msg.ErrTokenExpired
	}

	if err := s.blacklist.Add(ctx, tokenString, expiration); err != nil {
		return err_msg.ErrBlacklistAdd
	}

	return nil
}
