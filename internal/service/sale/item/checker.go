package services

import (
	"context"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *saleItemService) ItemExists(ctx context.Context, id int64) (bool, error) {
	if id <= 0 {
		return false, errMsg.ErrZeroID
	}

	return s.repo.ItemExists(ctx, id)
}
