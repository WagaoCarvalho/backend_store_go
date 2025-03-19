package services

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/WagaoCarvalho/backend_store_go/internal/repositories"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid int64) (models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *userService) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	user, err := s.repo.GetUserById(ctx, uid)
	if err != nil {
		return models.User{}, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return user, nil
}
