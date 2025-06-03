package auth

import (
	"context"
	"errors"
	"time"

	modelsLogin "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

var (
	ErrInvalidEmailFormat = errors.New("formato de email inválido")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserFetchFailed    = errors.New("erro ao buscar usuário")
	ErrTokenGeneration    = errors.New("erro ao gerar token de acesso")
	ErrAccountDisabled    = errors.New("conta desativada")
)

type LoginService interface {
	Login(ctx context.Context, credentials modelsLogin.LoginCredentials) (string, error)
}

type TokenGenerator interface {
	Generate(uid int64, email string) (string, error)
}

type loginService struct {
	userRepo   repositories.UserRepository
	jwtManager TokenGenerator
}

func NewLoginService(repo repositories.UserRepository, jwt TokenGenerator) *loginService {
	return &loginService{
		userRepo:   repo,
		jwtManager: jwt,
	}
}

func (s *loginService) Login(ctx context.Context, credentials modelsLogin.LoginCredentials) (string, error) {
	if !utils.IsValidEmail(credentials.Email) {
		return "", ErrInvalidEmailFormat
	}

	user, err := s.userRepo.GetByEmail(ctx, credentials.Email)
	if err != nil {
		time.Sleep(time.Second) // mitigação timing attack
		return "", ErrInvalidCredentials
	}

	if err := ComparePassword(user.Password, credentials.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.Status {
		return "", ErrAccountDisabled
	}

	token, err := s.jwtManager.Generate(user.UID, user.Email)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return token, nil
}
