package services

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userService) UserExists(ctx context.Context, userID int64) (bool, error) {
	if userID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return exists, nil
}
