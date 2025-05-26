package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
)

var (
	ErrInvalidAddressData = errors.New("dados do endereço inválidos")
	ErrAddressIDRequired  = errors.New("ID do endereço é obrigatório")
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
		return models.Address{}, ErrInvalidAddressData
	}
	return s.repo.Create(ctx, address)
}

func (s *addressService) GetByID(ctx context.Context, id int) (models.Address, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *addressService) Update(ctx context.Context, address models.Address) error {
	if address.ID == nil {
		return ErrAddressIDRequired
	}
	if address.Version == 0 {
		return errors.New("versão obrigatória para atualização")
	}

	err := s.repo.Update(ctx, address)
	if err != nil {
		if errors.Is(err, repositories.ErrVersionConflict) {
			return fmt.Errorf("conflito de versão: %w", err)
		}
		return fmt.Errorf("erro ao atualizar endereço: %w", err)
	}

	return nil
}

func (s *addressService) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return ErrAddressIDRequired
	}
	return s.repo.Delete(ctx, id)
}
