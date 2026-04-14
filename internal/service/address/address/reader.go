package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *addressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	return s.addressRepo.GetByID(ctx, id)
}
