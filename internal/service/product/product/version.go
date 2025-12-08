package services

import (
	"context"
	"errors"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productService) GetVersionByID(ctx context.Context, pid int64) (int64, error) {

	if pid <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repo.GetVersionByID(ctx, pid)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return 0, errMsg.ErrNotFound
		}

		return 0, fmt.Errorf("%w: %v", errMsg.ErrVersionConflict, err)
	}

	return version, nil
}
