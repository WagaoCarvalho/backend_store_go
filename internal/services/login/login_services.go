package services

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/user"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	Login(ctx context.Context, credentials models.LoginCredentials) (string, error)
}

type loginService struct {
	userRepo repositories.UserRepository
}

func NewLoginService(userRepo repositories.UserRepository) LoginService {
	return &loginService{userRepo: userRepo}
}

func (s *loginService) Login(ctx context.Context, credentials models.LoginCredentials) (string, error) {

	user, err := s.userRepo.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return "", fmt.Errorf("credenciais inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return "", fmt.Errorf("credenciais inválidas")
	}

	token, err := utils.GenerateJWT(user.UID, user.Email)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar token")
	}

	return token, nil
}
