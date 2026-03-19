package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

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
