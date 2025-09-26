package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoClient "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type AddressService interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetByUserID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetBySupplierID(ctx context.Context, clientID int64) ([]*models.Address, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
}

type addressService struct {
	repoAddress  repoAddress.AddressRepository
	repoClient   repoClient.ClientRepository
	repoUser     repoUser.UserRepository
	repoSupplier repoSupplier.SupplierRepository
}

func NewAddressService(
	repoAddress repoAddress.AddressRepository,
	repoClient repoClient.ClientRepository,
	repoUser repoUser.UserRepository,
	repoSupplier repoSupplier.SupplierRepository,
) AddressService {
	return &addressService{
		repoAddress:  repoAddress,
		repoClient:   repoClient,
		repoUser:     repoUser,
		repoSupplier: repoSupplier,
	}
}

func (s *addressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {

	if err := address.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	createdAddress, err := s.repoAddress.Create(ctx, address)
	if err != nil {
		return nil, err
	}

	return createdAddress, nil
}

func (s *addressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	if id <= 0 {
		return nil, errMsg.ErrIDZero
	}

	addressModel, err := s.repoAddress.GetByID(ctx, id)
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
		return nil, errMsg.ErrIDZero
	}

	addresses, err := s.repoAddress.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.repoUser.UserExists(ctx, userID)
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
		return nil, errMsg.ErrIDZero
	}

	addresses, err := s.repoAddress.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.repoClient.ClientExists(ctx, clientID)
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
		return nil, errMsg.ErrIDZero
	}

	addresses, err := s.repoAddress.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(addresses) == 0 {
		exists, err := s.repoSupplier.SupplierExists(ctx, supplierID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		if !exists {
			return nil, errMsg.ErrNotFound
		}
	}

	return addresses, nil
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {
	if address.ID <= 0 {
		return errMsg.ErrIDZero
	}
	if err := address.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if err := s.repoAddress.Update(ctx, address); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *addressService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrIDZero
	}

	if err := s.repoAddress.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (s *addressService) Disable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrIDZero
	}

	if err := s.repoAddress.Disable(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (s *addressService) Enable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrIDZero
	}

	if err := s.repoAddress.Enable(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}
