package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/auth"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repositories_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

type UserService interface {
	Create(ctx context.Context, user *models_user.User) (*models_user.User, error)
	GetAll(ctx context.Context) ([]*models_user.User, error)
	GetByID(ctx context.Context, uid int64) (*models_user.User, error)
	GetVersionByID(ctx context.Context, uid int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models_user.User, error)
	Delete(ctx context.Context, uid int64) error
	Update(ctx context.Context, user *models_user.User) (*models_user.User, error)
}

type userService struct {
	repo   repositories_user.UserRepository
	logger *logger.LoggerAdapter
	hasher auth.PasswordHasher
}

func NewUserService(repo repositories_user.UserRepository, logger *logger.LoggerAdapter, hasher auth.PasswordHasher) *userService {
	return &userService{
		repo:   repo,
		logger: logger,
		hasher: hasher,
	}
}

func (s *userService) Create(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	s.logger.Info(ctx, "[userService] - Iniciando criação de usuário", map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	})

	if !utils_validators.IsValidEmail(user.Email) {
		s.logger.Error(ctx, ErrInvalidEmail, "[userService] - Email inválido", map[string]interface{}{
			"email": user.Email,
		})
		return nil, ErrInvalidEmail
	}

	if user.Password != "" {
		hashed, err := s.hasher.Hash(user.Password)
		if err != nil {
			s.logger.Error(ctx, err, "[userService] - Erro ao hashear senha", map[string]interface{}{
				"email": user.Email,
			})
			return nil, fmt.Errorf("erro ao hashear senha: %w", err)
		}
		user.Password = hashed
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error(ctx, err, "[userService] - Erro ao criar usuário no repositório", map[string]interface{}{
			"email": user.Email,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	if createdUser == nil {
		s.logger.Error(ctx, nil, "[userService] - Usuário criado é nulo", map[string]interface{}{
			"email": user.Email,
		})
		return nil, fmt.Errorf("usuário criado é nulo")
	}

	s.logger.Info(ctx, "[userService] - Usuário criado com sucesso", map[string]interface{}{
		"user_id":  createdUser.UID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
	})

	return createdUser, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*models_user.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetByID(ctx context.Context, uid int64) (*models_user.User, error) {
	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}
	return user, nil
}

func (s *userService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	version, err := s.repo.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repositories_user.ErrUserNotFound) {
			return 0, repositories_user.ErrUserNotFound
		}
		return 0, fmt.Errorf("user: erro ao obter versão: %w", err)
	}
	return version, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}
	return user, nil
}

func (s *userService) Update(
	ctx context.Context,
	user *models_user.User,
) (*models_user.User, error) {
	if !utils_validators.IsValidEmail(user.Email) {
		return nil, ErrInvalidEmail
	}

	if user.Version <= 0 {
		return nil, ErrInvalidVersion
	}

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, repositories_user.ErrUserNotFound) {
			return nil, repositories_user.ErrUserNotFound
		}
		if errors.Is(err, repositories_user.ErrVersionConflict) {
			return nil, repositories_user.ErrVersionConflict
		}
		return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	err := s.repo.Delete(ctx, uid)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}
