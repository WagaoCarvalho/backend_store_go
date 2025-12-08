package services

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *clientService) GetVersionByID(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return version, nil
}
