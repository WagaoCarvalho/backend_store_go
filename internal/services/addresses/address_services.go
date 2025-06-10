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
	ErrVersionRequired    = errors.New("address: versão obrigatória para atualização")
	ErrUpdateAddress      = errors.New("address: erro ao atualizar")
	ErrVersionConflict    = errors.New("address: conflito de versão")
	ErrInvalidID          = errors.New("address: id inválido")
	ErrAddressNotFound    = errors.New("address: endereço não encontrado")
)

type AddressService interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetVersionByID(ctx context.Context, id int64) (int, error)
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

func (s *addressService) GetVersionByID(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return 0, ErrInvalidID
	}

	version, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrAddressNotFound) {
			return 0, ErrAddressNotFound
		}
		return 0, fmt.Errorf("failed to get address version: %w", err)
	}

	return version, nil
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {

	if err := address.Validate(); err != nil {
		return err
	}

	if address.ID == 0 {
		return ErrAddressIDRequired
	}
	if address.Version == 0 {
		return ErrVersionRequired
	}

	err := s.repo.Update(ctx, address)
	if err != nil {
		if errors.Is(err, repositories.ErrVersionConflict) {
			return fmt.Errorf("%w: %v", ErrVersionConflict, err)
		}
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
