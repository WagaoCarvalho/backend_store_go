package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	val_contact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
)

func (s *userService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	hashed, err := s.hasher.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}
	user.Password = hashed

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	if createdUser == nil {
		return nil, fmt.Errorf("usuário criado é nulo")
	}

	return createdUser, nil
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	if !val_contact.IsValidEmail(user.Email) {
		return errMsg.ErrInvalidData
	}

	if user.Version <= 0 {
		return errMsg.ErrVersionConflict
	}

	err := s.repo.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			return errMsg.ErrNotFound
		case errors.Is(err, errMsg.ErrVersionConflict):
			return errMsg.ErrVersionConflict
		default:
			return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {

	if uid <= 0 {
		return errMsg.ErrZeroID
	}
	return s.repo.Delete(ctx, uid)
}
