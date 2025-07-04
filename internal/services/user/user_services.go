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
	s.logger.Info(ctx, "[userService] - Iniciando recuperação de todos os usuários", nil)

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, err, "[userService] - Erro ao recuperar usuários do repositório", nil)
		return nil, fmt.Errorf("%w: %v", ErrGetAllUsers, err)
	}

	s.logger.Info(ctx, "[userService] - Recuperação de usuários concluída com sucesso", map[string]interface{}{
		"quantidade": len(users),
	})

	return users, nil
}

func (s *userService) GetByID(ctx context.Context, uid int64) (*models_user.User, error) {
	s.logger.Info(ctx, "[userService] - Iniciando recuperação de usuário por ID", map[string]interface{}{
		"user_id": uid,
	})

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, "[userService] - Erro ao recuperar usuário no repositório", map[string]interface{}{
			"user_id": uid,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	s.logger.Info(ctx, "[userService] - Usuário recuperado com sucesso", map[string]interface{}{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}

func (s *userService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	s.logger.Info(ctx, "[userService] - Iniciando recuperação de versão do usuário", map[string]interface{}{
		"user_id": uid,
	})

	version, err := s.repo.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repositories_user.ErrUserNotFound) {
			s.logger.Error(ctx, err, "[userService] - Usuário não encontrado ao buscar versão", map[string]interface{}{
				"user_id": uid,
			})
			return 0, repositories_user.ErrUserNotFound
		}

		s.logger.Error(ctx, err, "[userService] - Erro ao recuperar versão do usuário", map[string]interface{}{
			"user_id": uid,
		})
		return 0, fmt.Errorf("user: erro ao obter versão: %w", err)
	}

	s.logger.Info(ctx, "[userService] - Versão do usuário recuperada com sucesso", map[string]interface{}{
		"user_id": uid,
		"version": version,
	})

	return version, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	s.logger.Info(ctx, "[userService] - Iniciando recuperação de usuário por e-mail", map[string]interface{}{
		"email": email,
	})

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error(ctx, err, "[userService] - Erro ao recuperar usuário no repositório por e-mail", map[string]interface{}{
			"email": email,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	s.logger.Info(ctx, "[userService] - Usuário recuperado com sucesso por e-mail", map[string]interface{}{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}

func (s *userService) Update(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	s.logger.Info(ctx, "[userService] - Iniciando atualização de usuário", map[string]interface{}{
		"user_id":  user.UID,
		"email":    user.Email,
		"version":  user.Version,
		"username": user.Username,
	})

	if !utils_validators.IsValidEmail(user.Email) {
		s.logger.Error(ctx, ErrInvalidEmail, "[userService] - Email inválido", map[string]interface{}{
			"email": user.Email,
		})
		return nil, ErrInvalidEmail
	}

	if user.Version <= 0 {
		s.logger.Error(ctx, ErrInvalidVersion, "[userService] - Versão inválida para atualização", map[string]interface{}{
			"user_id": user.UID,
		})
		return nil, ErrInvalidVersion
	}

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repositories_user.ErrUserNotFound):
			s.logger.Error(ctx, err, "[userService] - Usuário não encontrado para atualização", map[string]interface{}{
				"user_id": user.UID,
			})
			return nil, repositories_user.ErrUserNotFound

		case errors.Is(err, repositories_user.ErrVersionConflict):
			s.logger.Error(ctx, err, "[userService] - Conflito de versão na atualização do usuário", map[string]interface{}{
				"user_id": user.UID,
				"version": user.Version,
			})
			return nil, repositories_user.ErrVersionConflict

		default:
			s.logger.Error(ctx, err, "[userService] - Erro inesperado ao atualizar usuário", map[string]interface{}{
				"user_id": user.UID,
			})
			return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
		}
	}

	s.logger.Info(ctx, "[userService] - Usuário atualizado com sucesso", map[string]interface{}{
		"user_id":  updatedUser.UID,
		"email":    updatedUser.Email,
		"username": updatedUser.Username,
		"version":  updatedUser.Version,
	})

	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	s.logger.Info(ctx, "[userService] - Iniciando exclusão de usuário", map[string]interface{}{
		"user_id": uid,
	})

	err := s.repo.Delete(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, "[userService] - Erro ao deletar usuário", map[string]interface{}{
			"user_id": uid,
		})
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	s.logger.Info(ctx, "[userService] - Usuário deletado com sucesso", map[string]interface{}{
		"user_id": uid,
	})

	return nil
}
