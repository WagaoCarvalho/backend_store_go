package auth

import (
	"context"
	"errors"
	"time"

	pass "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	logger "github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models_login "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

var (
	ErrInvalidEmailFormat = errors.New("formato de email inválido")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserFetchFailed    = errors.New("erro ao buscar usuário")
	ErrTokenGeneration    = errors.New("erro ao gerar token de acesso")
	ErrAccountDisabled    = errors.New("conta desativada")
)

type LoginService interface {
	Login(ctx context.Context, credentials models_login.LoginCredentials) (string, error)
}

type TokenGenerator interface {
	Generate(uid int64, email string) (string, error)
}

type loginService struct {
	userRepo   repositories.UserRepository
	logger     *logger.LoggerAdapter
	jwtManager TokenGenerator
	hasher     pass.PasswordHasher
}

func NewLoginService(repo repositories.UserRepository, logger *logger.LoggerAdapter, jwt TokenGenerator, hasher pass.PasswordHasher) *loginService {
	return &loginService{
		userRepo:   repo,
		logger:     logger,
		jwtManager: jwt,
		hasher:     hasher,
	}
}

func (s *loginService) Login(ctx context.Context, credentials models_login.LoginCredentials) (string, error) {
	ref := "[loginService - Login] - "
	s.logger.Info(ctx, ref+"iniciando login", map[string]any{
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
		s.logger.Warn(ctx, ref+"usuário não encontrado ou erro ao buscar", map[string]any{
			"email": credentials.Email,
		})
		return "", ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.Password, credentials.Password); err != nil {
		s.logger.Warn(ctx, ref+"senha inválida", map[string]any{
			"user_id": user.UID,
			"email":   credentials.Email,
		})
		return "", ErrInvalidCredentials
	}

	if !user.Status {
		s.logger.Warn(ctx, ref+"conta desativada", map[string]any{
			"user_id": user.UID,
			"email":   user.Email,
		})
		return "", ErrAccountDisabled
	}

	token, err := s.jwtManager.Generate(user.UID, user.Email)
	if err != nil {
		s.logger.Error(ctx, err, ref+"erro ao gerar token", map[string]any{
			"user_id": user.UID,
			"email":   user.Email,
		})
		return "", ErrTokenGeneration
	}

	s.logger.Info(ctx, ref+"login realizado com sucesso", map[string]any{
		"user_id": user.UID,
		"email":   user.Email,
	})

	return token, nil
}
