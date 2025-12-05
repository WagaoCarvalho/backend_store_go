package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	hashed, err := s.hasher.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: erro ao hashear senha", errMsg.ErrInternal)
	}
	user.Password = hashed

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		// CORREÇÃO: Envolver o erro com ErrCreate
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	if createdUser == nil {
		return nil, fmt.Errorf("%w: usuário criado é nulo", errMsg.ErrInternal)
	}

	return createdUser, nil
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return fmt.Errorf("%w: usuário não pode ser nulo", errMsg.ErrInvalidData)
	}

	if user.UID <= 0 {
		return fmt.Errorf("%w: ID do usuário inválido", errMsg.ErrZeroID)
	}

	if user.Version < 0 {
		return fmt.Errorf("%w: versão não pode ser negativa", errMsg.ErrInvalidData)
	}

	if err := user.ValidateForUpdate(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if user.Password != "" {

		if !strings.HasPrefix(user.Password, "$2a$") {
			hashed, err := s.hasher.Hash(user.Password)
			if err != nil {
				return fmt.Errorf("%w: erro ao processar senha", errMsg.ErrInternal)
			}
			user.Password = hashed

		}
	}

	err := s.repo.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			return fmt.Errorf("%w: usuário não encontrado", errMsg.ErrNotFound)
		case errors.Is(err, errMsg.ErrVersionConflict):
			return fmt.Errorf("%w: versão conflitante, dados desatualizados", errMsg.ErrVersionConflict)
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
