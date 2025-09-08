package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	val_contact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
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
	repoUser repo.UserRepository
	hasher   auth.PasswordHasher
}

func NewUserService(repoUser repo.UserRepository, hasher auth.PasswordHasher) UserService {
	return &userService{
		repoUser: repoUser,
		hasher:   hasher,
	}
}

func (s *userService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if !val_contact.IsValidEmail(user.Email) {
		return nil, err_msg.ErrInvalidData
	}

	if user.Password != "" {
		hashed, err := s.hasher.Hash(user.Password)
		if err != nil {
			return nil, fmt.Errorf("erro ao hashear senha: %w", err)
		}
		user.Password = hashed
	}

	createdUser, err := s.repoUser.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
	}

	if createdUser == nil {
		return nil, fmt.Errorf("usuário criado é nulo")
	}

	return createdUser, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := s.repoUser.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return users, nil
}

func (s *userService) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	if uid <= 0 {
		return nil, errors.New("ID inválido")
	}

	user, err := s.repoUser.GetByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return user, nil
}

func (s *userService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	version, err := s.repoUser.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return 0, err_msg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", err_msg.ErrVersionConflict, err)
	}
	return version, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}
	return user, nil
}

func (s *userService) GetByName(ctx context.Context, name string) ([]*models.User, error) {
	users, err := s.repoUser.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}
	return users, nil
}

func (s *userService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if !val_contact.IsValidEmail(user.Email) {
		return nil, err_msg.ErrInvalidData
	}

	if user.Version <= 0 {
		return nil, err_msg.ErrVersionConflict
	}

	updatedUser, err := s.repoUser.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, err_msg.ErrNotFound):
			return nil, err_msg.ErrNotFound
		case errors.Is(err, err_msg.ErrVersionConflict):
			return nil, err_msg.ErrVersionConflict
		default:
			return nil, fmt.Errorf("%w: %v", err_msg.ErrUpdate, err)
		}
	}

	return updatedUser, nil
}

func (s *userService) Disable(ctx context.Context, uid int64) error {
	return s.repoUser.Disable(ctx, uid)
}

func (s *userService) Enable(ctx context.Context, uid int64) error {
	err := s.repoUser.Enable(ctx, uid)
	if errors.Is(err, err_msg.ErrNotFound) {
		return err
	}
	return err
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	return s.repoUser.Delete(ctx, uid)
}
