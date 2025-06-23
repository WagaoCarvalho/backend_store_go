package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
)

var (
	ErrInvalidAddressData = errors.New("address: dados do endereço inválidos")
	ErrAddressIDRequired  = errors.New("address: ID do endereço é obrigatório")
	ErrUpdateAddress      = errors.New("address: erro ao atualizar")
	ErrInvalidID          = errors.New("address: id inválido")
	ErrAddressNotFound    = errors.New("address: endereço não encontrado")
)

type AddressService interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
}

type addressService struct {
	repo repositories.AddressRepository
}

func NewAddressService(repo repositories.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {

	if err := address.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, address)
}

func (s *addressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	if id == 0 {
		return nil, ErrAddressIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *addressService) GetByUserID(ctx context.Context, id int64) ([]*models.Address, error) {
	if id == 0 {
		return nil, ErrAddressIDRequired
	}
	return s.repo.GetByUserID(ctx, id)
}

func (s *addressService) GetByClientID(ctx context.Context, id int64) ([]*models.Address, error) {
	if id == 0 {
		return nil, ErrAddressIDRequired
	}
	return s.repo.GetByClientID(ctx, id)
}

func (s *addressService) GetBySupplierID(ctx context.Context, id int64) ([]*models.Address, error) {
	if id == 0 {
		return nil, ErrAddressIDRequired
	}
	return s.repo.GetBySupplierID(ctx, id)
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {
	if err := address.Validate(); err != nil {
		return err
	}

	if address.ID == 0 {
		return ErrAddressIDRequired
	}

	err := s.repo.Update(ctx, address)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateAddress, err)
	}

	return nil
}

func (s *addressService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrAddressIDRequired
	}

	return s.repo.Delete(ctx, id)
}
