package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"golang.org/x/crypto/bcrypt"
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
	// Validação básica
	if !utils.IsValidEmail(credentials.Email) {
		return "", fmt.Errorf("formato de email inválido")
	}
	// if len(credentials.Password) < 8 {
	// 	return "", fmt.Errorf("a senha deve ter pelo menos 8 caracteres")
	// }

	user, err := s.userRepo.GetByEmail(ctx, credentials.Email)
	if err != nil {
		log.Printf("Erro ao buscar usuário: %v", err)
		if errors.Is(err, repositories.ErrUserNotFound) {
			// Delay para prevenir timing attacks
			time.Sleep(time.Second)
			return "", fmt.Errorf("credenciais inválidas")
		}
		return "", fmt.Errorf("erro ao buscar usuário")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return "", fmt.Errorf("credenciais inválidas")
	}

	if !user.Status {
		return "", fmt.Errorf("conta desativada")
	}

	// 👇 Aqui estava o problema
	token, err := s.generateJWT(user.UID, user.Email)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar token de acesso")
	}

	return token, nil
}
