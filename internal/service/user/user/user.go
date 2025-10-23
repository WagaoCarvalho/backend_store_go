package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	val_contact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type User interface {
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

type user struct {
	repoUser repo.User
	hasher   auth.PasswordHasher
}

func NewUser(repoUser repo.User, hasher auth.PasswordHasher) User {
	return &user{
		repoUser: repoUser,
		hasher:   hasher,
	}
}

func (s *user) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	hashed, err := s.hasher.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}
	user.Password = hashed

	createdUser, err := s.repoUser.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	if createdUser == nil {
		return nil, fmt.Errorf("usuário criado é nulo")
	}

	return createdUser, nil
}

func (s *user) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := s.repoUser.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return users, nil
}

func (s *user) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	if uid <= 0 {
		return nil, errMsg.ErrZeroID
	}

	user, err := s.repoUser.GetByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return user, nil
}

func (s *user) GetVersionByID(ctx context.Context, uid int64) (int64, error) {

	if uid <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repoUser.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return 0, errMsg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", errMsg.ErrVersionConflict, err)
	}
	return version, nil
}

func (s *user) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email inválido")
	}

	user, err := s.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return user, nil
}

func (s *user) GetByName(ctx context.Context, name string) ([]*models.User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("nome inválido")
	}
	users, err := s.repoUser.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return users, nil
}

func (s *user) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if !val_contact.IsValidEmail(user.Email) {
		return nil, errMsg.ErrInvalidData
	}

	if user.Version <= 0 {
		return nil, errMsg.ErrVersionConflict
	}

	updatedUser, err := s.repoUser.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			return nil, errMsg.ErrNotFound
		case errors.Is(err, errMsg.ErrVersionConflict):
			return nil, errMsg.ErrVersionConflict
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
		}
	}

	return updatedUser, nil
}

func (s *user) Disable(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return errMsg.ErrZeroID
	}
	return s.repoUser.Disable(ctx, uid)
}

func (s *user) Enable(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return errMsg.ErrZeroID
	}
	err := s.repoUser.Enable(ctx, uid)
	if errors.Is(err, errMsg.ErrNotFound) {
		return err
	}
	return err
}

func (s *user) Delete(ctx context.Context, uid int64) error {

	if uid <= 0 {
		return errMsg.ErrZeroID
	}
	return s.repoUser.Delete(ctx, uid)
}
