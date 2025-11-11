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
	if userID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := s.addressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.userRepo.UserExists(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}

func (s *addressService) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
	if clientID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := s.addressRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.clientRepo.ClientExists(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}

func (s *addressService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	if supplierID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := s.addressRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.supplierRepo.SupplierExists(ctx, supplierID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}
