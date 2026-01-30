package services

import (
	"context"
	"errors"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *addressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	if address == nil {
		return nil, errMsg.ErrNilModel
	}

	if err := address.Validate(); err != nil {
		return nil, errors.Join(errMsg.ErrInvalidData, err)
	}

	created, err := s.addressRepo.Create(ctx, address)
	if err != nil {
		return nil, errors.Join(errMsg.ErrCreate, err)
	}

	return created, nil
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {
	if address == nil {
		return errMsg.ErrNilModel
	}
	if address.ID <= 0 {
		return errMsg.ErrZeroID
	}
	if err := address.Validate(); err != nil {
		return errors.Join(errMsg.ErrInvalidData, err)
	}

	if err := s.addressRepo.Update(ctx, address); err != nil {
		return errors.Join(errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *addressService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.addressRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return errors.Join(errMsg.ErrDelete, err)
	}

	return nil
}
