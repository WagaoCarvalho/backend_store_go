package auth

import (
	"context"
	"time"

	pass "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
	logger "github.com/WagaoCarvalho/backend_store_go/pkg/logger"
)

type LoginService interface {
	Login(ctx context.Context, credentials model.LoginCredentials) (string, error)
}

type TokenGenerator interface {
	Generate(uid int64, email string) (string, error)
}

type loginService struct {
	userRepo   repo.UserRepository
	logger     *logger.LoggerAdapter
	jwtManager TokenGenerator
	hasher     pass.PasswordHasher
}

func NewLoginService(repo repo.UserRepository, logger *logger.LoggerAdapter, jwt TokenGenerator, hasher pass.PasswordHasher) *loginService {
	return &loginService{
		userRepo:   repo,
		logger:     logger,
		jwtManager: jwt,
		hasher:     hasher,
	}
}

func (s *loginService) Login(ctx context.Context, credentials model.LoginCredentials) (string, error) {
	const ref = "[loginService - Login] - "

	s.logger.Info(ctx, ref+logger.LogLoginInit, map[string]any{
		"email": credentials.Email,
	})

	if !utils_validators.IsValidEmail(credentials.Email) {
		s.logger.Error(ctx, ErrInvalidEmailFormat, ref+logger.LogEmailInvalid, map[string]any{
			"email": credentials.Email,
		})
		return "", ErrInvalidEmailFormat
	}

	user, err := s.userRepo.GetByEmail(ctx, credentials.Email)
	if err != nil {
		time.Sleep(time.Second) // mitigação timing attack
		s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"email": credentials.Email,
		})
		return "", ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.Password, credentials.Password); err != nil {
		s.logger.Warn(ctx, ref+logger.LogPasswordInvalid, map[string]any{
			"user_id": user.UID,
			"email":   credentials.Email,
		})
		return "", ErrInvalidCredentials
	}

	if !user.Status {
		s.logger.Warn(ctx, ref+logger.LogAccountDisabled, map[string]any{
			"user_id": user.UID,
			"email":   user.Email,
		})
		return "", ErrAccountDisabled
	}

	token, err := s.jwtManager.Generate(user.UID, user.Email)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogTokenGenerationError, map[string]any{
			"user_id": user.UID,
			"email":   user.Email,
		})
		return "", ErrTokenGeneration
	}

	s.logger.Info(ctx, ref+logger.LogLoginSuccess, map[string]any{
		"user_id": user.UID,
		"email":   user.Email,
	})

	return token, nil
}
