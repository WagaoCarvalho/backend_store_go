package services

import (
	"context"
	"errors"
	"fmt"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/users"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, uid int64) (*models.User, error)
	GetVersionByID(ctx context.Context, uid int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByName(ctx context.Context, name string) ([]*models.User, error)
	Delete(ctx context.Context, uid int64) error
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
	Update(ctx context.Context, user *models.User) (*models.User, error)
}

type userService struct {
	repo_user repo.UserRepository
	logger    *logger.LoggerAdapter
	hasher    auth.PasswordHasher
}

func NewUserService(repo_user repo.UserRepository, logger *logger.LoggerAdapter, hasher auth.PasswordHasher) UserService {
	return &userService{
		repo_user: repo_user,
		logger:    logger,
		hasher:    hasher,
	}
}

func (s *userService) Create(ctx context.Context, user *models.User) (*models.User, error) {
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

	createdUser, err := s.repo_user.Create(ctx, user)
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

func (s *userService) GetAll(ctx context.Context) ([]*models.User, error) {
	ref := "[userService - GetAll] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, nil)

	users, err := s.repo_user.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrGetAllUsers, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"quantidade": len(users),
	})

	return users, nil
}

func (s *userService) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	ref := "[userService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": uid,
	})

	if uid <= 0 {
		s.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"product_id": uid,
		})
		return nil, errors.New("ID inválido")
	}

	user, err := s.repo_user.GetByID(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": uid,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}

func (s *userService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	ref := "[userService - GetVersionByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": uid,
	})

	version, err := s.repo_user.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			s.logger.Error(ctx, err, ref+logger.LogNotFound, map[string]any{
				"user_id": uid,
			})
			return 0, repo.ErrUserNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": uid,
		})
		return 0, fmt.Errorf("%w: %v", ErrInvalidVersion, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": uid,
		"version": version,
	})

	return version, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ref := "[userService - GetByEmail] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"email": email,
	})

	user, err := s.repo_user.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"email": email,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}

func (s *userService) GetByName(ctx context.Context, name string) ([]*models.User, error) {
	ref := "[userService - GetByName] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"username_partial": name,
	})

	users, err := s.repo_user.GetByName(ctx, name)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"username_partial": name,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(users),
	})

	return users, nil
}

func (s *userService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ref := "[userService - Update] - "

	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"user_id":  user.UID,
		"email":    user.Email,
		"version":  user.Version,
		"username": user.Username,
	})

	if !utils_validators.IsValidEmail(user.Email) {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"email": user.Email,
		})
		return nil, ErrInvalidEmail
	}

	if user.Version <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": user.UID,
			"version": user.Version,
		})
		return nil, ErrInvalidVersion
	}

	updatedUser, err := s.repo_user.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": user.UID,
			})
			return nil, repo.ErrUserNotFound

		case errors.Is(err, repo.ErrVersionConflict):
			s.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"user_id": user.UID,
				"version": user.Version,
			})
			return nil, repo.ErrVersionConflict

		default:
			s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"user_id": user.UID,
			})
			return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id":  updatedUser.UID,
		"email":    updatedUser.Email,
		"username": updatedUser.Username,
		"version":  updatedUser.Version,
	})

	return updatedUser, nil
}

func (s *userService) Disable(ctx context.Context, uid int64) error {
	ref := "[userService - Disable] - "
	s.logger.Info(ctx, ref+logger.LogDisableInit, map[string]any{
		"user_id": uid,
	})

	err := s.repo_user.Disable(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDisableError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrDisableUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogDisableSuccess, map[string]any{
		"user_id": uid,
	})

	return nil
}

func (s *userService) Enable(ctx context.Context, uid int64) error {
	ref := "[userService - Enable] - "

	s.logger.Info(ctx, ref+logger.LogEnableInit, map[string]any{
		"user_id": uid,
	})

	err := s.repo_user.Enable(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": uid,
			})
			return err
		default:
			s.logger.Error(ctx, err, ref+logger.LogEnableError, map[string]any{
				"user_id": uid,
			})
			return err
		}
	}

	s.logger.Info(ctx, ref+logger.LogEnableSuccess, map[string]any{
		"user_id": uid,
		"status":  true,
	})

	return nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	ref := "[userService - Delete] - "

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id": uid,
	})

	err := s.repo_user.Delete(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrDeleteUser, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id": uid,
	})

	return nil
}
