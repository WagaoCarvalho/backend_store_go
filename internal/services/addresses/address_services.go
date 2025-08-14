package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
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
	repo   repo.AddressRepository
	logger *logger.LoggerAdapter
}

func NewAddressService(repo repo.AddressRepository, logger *logger.LoggerAdapter) AddressService {
	return &addressService{
		repo:   repo,
		logger: logger,
	}
}

func (s *addressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	ref := "[addressService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
	})

	if err := address.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		return nil, err
	}

	createdAddress, err := s.repo.Create(ctx, address)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"street": address.Street,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"address_id": createdAddress.ID,
	})

	return createdAddress, nil
}

func (s *addressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	ref := "[addressService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"address_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"address_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	address, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrAddressNotFound) {
			s.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
				"address_id": id,
			})
			return nil, ErrAddressNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"address_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrCheckAddress, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"address_id": address.ID,
	})

	return address, nil
}

func (s *addressService) GetByUserID(ctx context.Context, id int64) ([]*models.Address, error) {
	ref := "[addressService - GetByUserID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	addresses, err := s.repo.GetByUserID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":         id,
		"total_addresses": len(addresses),
	})

	return addresses, nil
}

func (s *addressService) GetByClientID(ctx context.Context, id int64) ([]*models.Address, error) {
	ref := "[addressService - GetByClientID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"client_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"client_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	addresses, err := s.repo.GetByClientID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id":   id,
		"total_items": len(addresses),
	})

	return addresses, nil
}

func (s *addressService) GetBySupplierID(ctx context.Context, id int64) ([]*models.Address, error) {
	ref := "[addressService - GetBySupplierID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	addresses, err := s.repo.GetBySupplierID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": id,
		"total_items": len(addresses),
	})

	return addresses, nil
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {
	ref := "[addressService - Update] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"address_id":  address.ID,
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
	})

	if err := address.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		return err
	}

	if address.ID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"address_id": address.ID,
		})
		return ErrAddressIDRequired
	}

	err := s.repo.Update(ctx, address)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"address_id": address.ID,
		})
		return fmt.Errorf("%w: %v", ErrUpdateAddress, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"address_id": address.ID,
	})

	return nil
}
func (s *addressService) Delete(ctx context.Context, id int64) error {
	ref := "[addressService - Delete] - "
	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"address_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"address_id": id,
		})
		return ErrAddressIDRequired
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"address_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAddress, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"address_id": id,
	})
	return nil
}
