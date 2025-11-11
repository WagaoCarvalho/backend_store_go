package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *addressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {

	if err := address.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	createdAddress, err := s.address.Create(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return createdAddress, nil
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {
	if address.ID <= 0 {
		return errMsg.ErrZeroID
	}
	if err := address.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if err := s.address.Update(ctx, address); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *addressService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	_, err := s.address.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if err := s.address.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
