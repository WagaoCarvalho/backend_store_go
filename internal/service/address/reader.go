package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *address) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addressModel, err := s.address.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return addressModel, nil
}

func (s *address) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	if userID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := s.address.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.user.UserExists(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}

func (s *address) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
	if clientID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := s.address.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.client.ClientExists(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}

func (s *address) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	if supplierID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	addresses, err := s.address.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.supplier.SupplierExists(ctx, supplierID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}
