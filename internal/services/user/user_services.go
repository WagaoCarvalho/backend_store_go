package services

import (
	"context"
	"fmt"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type UserService interface {
	GetAll(ctx context.Context) ([]models_user.User, error)
	GetById(ctx context.Context, uid int64) (models_user.User, error)
	GetByEmail(ctx context.Context, email string) (models_user.User, error)
	Delete(ctx context.Context, uid int64) error
	Update(ctx context.Context, user models_user.User, contact *models_contact.Contact) (models_user.User, error)
	Create(
		ctx context.Context,
		user models_user.User,
		categoryID int64,
		address models_address.Address,
		contact models_contact.Contact,
	) (models_user.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll(ctx context.Context) ([]models_user.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetById(ctx context.Context, uid int64) (models_user.User, error) {
	user, err := s.repo.GetById(ctx, uid)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (models_user.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return user, nil
}

func (s *userService) Create(
	ctx context.Context,
	user models_user.User,
	categoryID int64,
	address models_address.Address,
	contact models_contact.Contact,
) (models_user.User, error) {

	if !utils.IsValidEmail(user.Email) {
		return models_user.User{}, fmt.Errorf("email inválido")
	}

	createdUser, err := s.repo.Create(ctx, user, categoryID, address, contact)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return createdUser, nil
}

func (s *userService) Update(ctx context.Context, user models_user.User, contact *models_contact.Contact) (models_user.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return models_user.User{}, fmt.Errorf("email inválido")
	}

	updatedUser, err := s.repo.Update(ctx, user, contact)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	if err := s.repo.Delete(ctx, uid); err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}
