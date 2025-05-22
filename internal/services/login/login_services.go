package services

import (
	"context"
	"errors"
	"log"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"golang.org/x/crypto/bcrypt"
)

// Erros personalizados
var (
	ErrInvalidEmailFormat = errors.New("formato de email inválido")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserFetchFailed    = errors.New("erro ao buscar usuário")
	ErrTokenGeneration    = errors.New("erro ao gerar token de acesso")
	ErrAccountDisabled    = errors.New("conta desativada")
)

type LoginService interface {
	Login(ctx context.Context, credentials models.LoginCredentials) (string, error)
}

type JWTGenerator func(uid int64, email string) (string, error)

type loginService struct {
	userRepo    repositories.UserRepository
	generateJWT JWTGenerator
}

func NewLoginService(repo repositories.UserRepository) *loginService {
	return &loginService{
		userRepo:    repo,
		generateJWT: utils.GenerateJWT,
	}
}

func NewLoginServiceWithJWT(repo repositories.UserRepository, jwtGen JWTGenerator) *loginService {
	return &loginService{
		userRepo:    repo,
		generateJWT: jwtGen,
	}
}

func (s *loginService) Login(ctx context.Context, credentials models.LoginCredentials) (string, error) {
	if !utils.IsValidEmail(credentials.Email) {
		return "", ErrInvalidEmailFormat
	}
	// Validação de senha desativada, mas deixada como sugestão:
	// if len(credentials.Password) < 8 {
	// 	return "", ErrWeakPassword
	// }

	user, err := s.userRepo.GetByEmail(ctx, credentials.Email)
	if err != nil {
		log.Printf("Erro ao buscar usuário: %v", err)
		if errors.Is(err, repositories.ErrUserNotFound) {
			time.Sleep(time.Second) // Mitigação de timing attack
			return "", ErrInvalidCredentials
		}
		return "", ErrUserFetchFailed
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.Status {
		return "", ErrAccountDisabled
	}

	token, err := s.generateJWT(user.UID, user.Email)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return token, nil
}
