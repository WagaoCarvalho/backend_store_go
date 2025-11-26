package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *addressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addressModel, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return addressModel, nil
}

func (s *addressService) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	return getAddressesByEntity(
		ctx,
		userID,
		s.addressRepo.GetByUserID,
		s.userRepo.UserExists,
	)
}

func (s *addressService) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
	return getAddressesByEntity(
		ctx,
		clientID,
		s.addressRepo.GetByClientID,
		s.clientRepo.ClientExists,
	)
}

func (s *addressService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	return getAddressesByEntity(
		ctx,
		supplierID,
		s.addressRepo.GetBySupplierID,
		s.supplierRepo.SupplierExists,
	)
}

func getAddressesByEntity(
	ctx context.Context,
	id int64,
	findFn func(context.Context, int64) ([]*models.Address, error),
	existsFn func(context.Context, int64) (bool, error),
) ([]*models.Address, error) {

	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := findFn(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := existsFn(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}
