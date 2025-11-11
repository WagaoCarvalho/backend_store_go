package services

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userService) Disable(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return errMsg.ErrZeroID
	}

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if !user.Status {
		return nil
	}

	if err := s.repo.Disable(ctx, uid); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (s *userService) Enable(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return errMsg.ErrZeroID
	}

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if user.Status {
		return nil
	}

	if err := s.repo.Enable(ctx, uid); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}
