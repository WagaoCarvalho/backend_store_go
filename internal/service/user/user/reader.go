package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userService) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return users, nil
}

func (s *userService) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	if uid <= 0 {
		return nil, errMsg.ErrZeroID
	}

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return user, nil
}

func (s *userService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {

	if uid <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repo.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return 0, errMsg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", errMsg.ErrVersionConflict, err)
	}
	return version, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email inválido")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return user, nil
}

func (s *userService) GetByName(ctx context.Context, name string) ([]*models.User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("nome inválido")
	}
	users, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return users, nil
}
