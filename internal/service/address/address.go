package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
)

type AddressService interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetByUserID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetBySupplierID(ctx context.Context, clientID int64) ([]*models.Address, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
}

type addressService struct {
	repo repo.AddressRepository
}

func NewAddressService(repo repo.AddressRepository) AddressService {
	return &addressService{
		repo: repo,
	}
}

func (s *addressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {

	if err := address.Validate(); err != nil {
		return nil, err
	}

	createdAddress, err := s.repo.Create(ctx, address)
	if err != nil {
		return nil, err
	}

	return createdAddress, nil
}

func (s *addressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	if id <= 0 {
		return nil, err_msg.ErrID
	}

	addressModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return nil, err_msg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return addressModel, nil
}

func (s *addressService) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	if userID <= 0 {
		return nil, err_msg.ErrID
	}

	addressModels, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return addressModels, nil
}

func (s *addressService) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
	if clientID <= 0 {
		return nil, err_msg.ErrID
	}

	addressModels, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return addressModels, nil
}

func (s *addressService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	if supplierID <= 0 {
		return nil, err_msg.ErrID
	}

	addressModels, err := s.repo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, err
	}

	return addressModels, nil
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {
	if err := address.Validate(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, address); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrUpdate, err)
	}

	return nil
}

func (s *addressService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return err_msg.ErrID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
