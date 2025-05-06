package services

import (
	"context"
	"errors"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
)

type AddressService interface {
	Create(ctx context.Context, address models.Address) (models.Address, error)
	GetByID(ctx context.Context, id int) (models.Address, error)
	Update(ctx context.Context, address models.Address) error
	Delete(ctx context.Context, id int) error
}

type addressService struct {
	repo repositories.AddressRepository
}

func NewAddressService(repo repositories.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) Create(ctx context.Context, address models.Address) (models.Address, error) {
	if address.Street == "" || address.City == "" || address.State == "" || address.PostalCode == "" {
		return models.Address{}, errors.New("dados do endereço inválidos")
	}

	return s.repo.Create(ctx, address)
}

func (s *addressService) GetByID(ctx context.Context, id int) (models.Address, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *addressService) Update(ctx context.Context, address models.Address) error {
	if address.ID == nil {
		return errors.New("ID do endereço é obrigatório")
	}

	return s.repo.Update(ctx, address)
}

func (s *addressService) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("ID do endereço é obrigatório")
	}

	return s.repo.Delete(ctx, id)
}
