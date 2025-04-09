package services

import (
	"context"
	"errors"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
)

type AddressService interface {
	CreateAddress(ctx context.Context, address models.Address) (models.Address, error)
	GetAddressByID(ctx context.Context, id int) (models.Address, error)
	UpdateAddress(ctx context.Context, address models.Address) error
	DeleteAddress(ctx context.Context, id int) error
}

type addressService struct {
	repo repositories.AddressRepository
}

func NewAddressService(repo repositories.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) CreateAddress(ctx context.Context, address models.Address) (models.Address, error) {
	if address.Street == "" || address.City == "" || address.State == "" || address.PostalCode == "" {
		return models.Address{}, errors.New("dados do endereço inválidos")
	}

	return s.repo.CreateAddress(ctx, address)
}

func (s *addressService) GetAddressByID(ctx context.Context, id int) (models.Address, error) {
	return s.repo.GetAddressByID(ctx, id)
}

func (s *addressService) UpdateAddress(ctx context.Context, address models.Address) error {
	if address.ID == 0 {
		return errors.New("ID do endereço é obrigatório")
	}

	return s.repo.UpdateAddress(ctx, address)
}

func (s *addressService) DeleteAddress(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("ID do endereço é obrigatório")
	}

	return s.repo.DeleteAddress(ctx, id)
}
