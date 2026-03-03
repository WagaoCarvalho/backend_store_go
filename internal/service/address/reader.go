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

func (s *addressService) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	return getAddressesByEntity(
		ctx,
		userID,
		s.addressRepo.GetByUserID,
	)
}

func (s *addressService) GetByClientCpfID(ctx context.Context, clientCpfID int64) ([]*models.Address, error) {
	return getAddressesByEntity(
		ctx,
		clientCpfID,
		s.addressRepo.GetByClientCpfID,
	)
}

func (s *addressService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	return getAddressesByEntity(
		ctx,
		supplierID,
		s.addressRepo.GetBySupplierID,
	)
}

func getAddressesByEntity(
	ctx context.Context,
	id int64,
	findFn func(context.Context, int64) ([]*models.Address, error),
) ([]*models.Address, error) {

	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := findFn(ctx, id)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}
