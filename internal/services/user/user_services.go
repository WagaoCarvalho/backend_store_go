package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid int64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User, categoryID int64) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	DeleteUserById(ctx context.Context, uid int64) error
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

func (s *userService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, user models.User, categoryID int64) (models.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return models.User{}, fmt.Errorf("email inválido")
	}

	createdUser, err := s.repo.CreateUser(ctx, user, categoryID)
	if err != nil {
		return models.User{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return createdUser, nil
}

func (s *userService) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return models.User{}, fmt.Errorf("email inválido")
	}

	updatedUser, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	return updatedUser, nil
}

func (s *userService) DeleteUserById(ctx context.Context, uid int64) error {
	if err := s.repo.DeleteUserById(ctx, uid); err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}
