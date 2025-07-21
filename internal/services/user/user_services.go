package services

import (
	"context"
	"errors"
	"fmt"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
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
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
	Update(ctx context.Context, user *models_user.User) (*models_user.User, error)
}

type userService struct {
	repo   repositories_user.UserRepository
	logger *logger.LoggerAdapter
	hasher auth.PasswordHasher
}

func NewUserService(repo repositories_user.UserRepository, logger *logger.LoggerAdapter, hasher auth.PasswordHasher) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
		hasher: hasher,
	}
}

func (s *userService) Create(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	ref := "[userService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	})

	if !utils_validators.IsValidEmail(user.Email) {
		s.logger.Error(ctx, ErrInvalidEmail, ref+logger.LogEmailInvalid, map[string]any{
			"email": user.Email,
		})
		return nil, ErrInvalidEmail
	}

	if user.Password != "" {
		hashed, err := s.hasher.Hash(user.Password)
		if err != nil {
			s.logger.Error(ctx, err, ref+logger.LogPasswordInvalid, map[string]any{
				"email": user.Email,
			})
			return nil, fmt.Errorf("erro ao hashear senha: %w", err)
		}
		user.Password = hashed
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"email": user.Email,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	if createdUser == nil {
		s.logger.Error(ctx, nil, ref+"usuário criado é nulo", map[string]any{
			"email": user.Email,
		})
		return nil, fmt.Errorf("usuário criado é nulo")
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdUser.UID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
	})

	return createdUser, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*models_user.User, error) {
	ref := "[userService - GetAll] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, nil)

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrGetAllUsers, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"quantidade": len(users),
	})

	return users, nil
}

func (s *userService) GetByID(ctx context.Context, uid int64) (*models_user.User, error) {
	ref := "[userService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]interface{}{
		"user_id": uid,
	})

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]interface{}{
			"user_id": uid,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]interface{}{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}

func (s *userService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	ref := "[userService - GetVersionByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]interface{}{
		"user_id": uid,
	})

	version, err := s.repo.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repositories_user.ErrUserNotFound) {
			s.logger.Error(ctx, err, ref+logger.LogNotFound, map[string]interface{}{
				"user_id": uid,
			})
			return 0, repositories_user.ErrUserNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]interface{}{
			"user_id": uid,
		})
		return 0, fmt.Errorf("%w: %v", ErrInvalidVersion, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]interface{}{
		"user_id": uid,
		"version": version,
	})

	return version, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	ref := "[userService - GetByEmail] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]interface{}{
		"email": email,
	})

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]interface{}{
			"email": email,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]interface{}{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}

func (s *userService) Update(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	ref := "[userService - Update] - "

	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]interface{}{
		"user_id":  user.UID,
		"email":    user.Email,
		"version":  user.Version,
		"username": user.Username,
	})

	if !utils_validators.IsValidEmail(user.Email) {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]interface{}{
			"email": user.Email,
		})
		return nil, ErrInvalidEmail
	}

	if user.Version <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]interface{}{
			"user_id": user.UID,
			"version": user.Version,
		})
		return nil, ErrInvalidVersion
	}

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repositories_user.ErrUserNotFound):
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]interface{}{
				"user_id": user.UID,
			})
			return nil, repositories_user.ErrUserNotFound

		case errors.Is(err, repositories_user.ErrVersionConflict):
			s.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]interface{}{
				"user_id": user.UID,
				"version": user.Version,
			})
			return nil, repositories_user.ErrVersionConflict

		default:
			s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]interface{}{
				"user_id": user.UID,
			})
			return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]interface{}{
		"user_id":  updatedUser.UID,
		"email":    updatedUser.Email,
		"username": updatedUser.Username,
		"version":  updatedUser.Version,
	})

	return updatedUser, nil
}

func (s *userService) Disable(ctx context.Context, uid int64) error {
	ref := "[userService - Disable] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"user_id": uid,
	})

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	user.Status = false

	_, err = s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrDisableUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id": uid,
		"status":  user.Status,
	})

	return nil
}

func (s *userService) Enable(ctx context.Context, uid int64) error {
	ref := "[userService - Enable] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"user_id": uid,
	})

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	user.Status = true

	_, err = s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrEnableUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id": uid,
		"status":  user.Status,
	})

	return nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	ref := "[userService - Delete] - "

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]interface{}{
		"user_id": uid,
	})

	err := s.repo.Delete(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]interface{}{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrDeleteUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]interface{}{
		"user_id": uid,
	})

	return nil
}
