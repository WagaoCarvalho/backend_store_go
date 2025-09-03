package services

import (
	"context"
	"errors"
	"fmt"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
)

type AddressService interface {
	Create(ctx context.Context, dto *dtoAddress.AddressDTO) (*dtoAddress.AddressDTO, error)
	GetByID(ctx context.Context, id int64) (*dtoAddress.AddressDTO, error)
	GetByUserID(ctx context.Context, userID int64) ([]*dtoAddress.AddressDTO, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*dtoAddress.AddressDTO, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*dtoAddress.AddressDTO, error)
	Update(ctx context.Context, addressDTO *dtoAddress.AddressDTO) error
	Delete(ctx context.Context, id int64) error
}

type addressService struct {
	repo   repo.AddressRepository
	logger *logger.LogAdapter
}

func NewAddressService(repo repo.AddressRepository, logger *logger.LogAdapter) AddressService {
	return &addressService{
		repo:   repo,
		logger: logger,
	}
}

func (s *addressService) Create(ctx context.Context, addressDTO *dtoAddress.AddressDTO) (*dtoAddress.AddressDTO, error) {
	ref := "[addressService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"user_id":     utils.Int64OrNil(addressDTO.UserID),
		"client_id":   utils.Int64OrNil(addressDTO.ClientID),
		"supplier_id": utils.Int64OrNil(addressDTO.SupplierID),
	})

	addressModel := dtoAddress.ToAddressModel(*addressDTO)

	if err := addressModel.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		return nil, err
	}

	createdAddress, err := s.repo.Create(ctx, addressModel)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"street": addressModel.Street,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"address_id": createdAddress.ID,
	})

	result := dtoAddress.ToAddressDTO(createdAddress)
	return &result, nil
}

func (s *addressService) GetByID(ctx context.Context, id int64) (*dtoAddress.AddressDTO, error) {
	ref := "[addressService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"address_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"address_id": id,
		})
		return nil, err_msg.ErrID
	}

	addressModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			s.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
				"address_id": id,
			})
			return nil, err_msg.ErrNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"address_id": id,
		})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	addressDTO := dtoAddress.ToAddressDTO(addressModel)

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"address_id": addressDTO.ID,
	})

	return &addressDTO, nil
}

func (s *addressService) GetByUserID(ctx context.Context, userID int64) ([]*dtoAddress.AddressDTO, error) {
	ref := "[addressService - GetByUserID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})

	if userID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return nil, err_msg.ErrID
	}

	addressModels, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, err
	}

	addressDTOs := make([]*dtoAddress.AddressDTO, len(addressModels))
	for i, addr := range addressModels {
		dto := dtoAddress.ToAddressDTO(addr)
		addressDTOs[i] = &dto

	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":         userID,
		"total_addresses": len(addressDTOs),
	})

	return addressDTOs, nil
}

func (s *addressService) GetByClientID(ctx context.Context, clientID int64) ([]*dtoAddress.AddressDTO, error) {
	ref := "[addressService - GetByClientID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"client_id": clientID,
	})

	if clientID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"client_id": clientID,
		})
		return nil, err_msg.ErrID
	}

	addressModels, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": clientID,
		})
		return nil, err
	}

	addressDTOs := make([]*dtoAddress.AddressDTO, len(addressModels))
	for i, addr := range addressModels {
		dto := dtoAddress.ToAddressDTO(addr)
		addressDTOs[i] = &dto
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id":   clientID,
		"total_items": len(addressDTOs),
	})

	return addressDTOs, nil
}

func (s *addressService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*dtoAddress.AddressDTO, error) {
	ref := "[addressService - GetBySupplierID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
	})

	if supplierID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, err_msg.ErrID
	}

	addressModels, err := s.repo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, err
	}

	addressDTOs := make([]*dtoAddress.AddressDTO, len(addressModels))
	for i, addr := range addressModels {
		dto := dtoAddress.ToAddressDTO(addr)
		addressDTOs[i] = &dto
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplierID,
		"total_items": len(addressDTOs),
	})

	return addressDTOs, nil
}

func (s *addressService) Update(ctx context.Context, addressDTO *dtoAddress.AddressDTO) error {
	ref := "[addressService - Update] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"address_id":  utils.Int64OrNil(addressDTO.ID),
		"user_id":     utils.Int64OrNil(addressDTO.UserID),
		"client_id":   utils.Int64OrNil(addressDTO.ClientID),
		"supplier_id": utils.Int64OrNil(addressDTO.SupplierID),
	})

	addressModel := dtoAddress.ToAddressModel(*addressDTO)

	if err := addressModel.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		return err
	}

	if addressModel.ID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"address_id": addressModel.ID,
		})
		return err_msg.ErrID
	}

	if err := s.repo.Update(ctx, addressModel); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"address_id": addressModel.ID,
		})
		return fmt.Errorf("%w: %v", err_msg.ErrUpdate, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"address_id": addressModel.ID,
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
		return err_msg.ErrID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"address_id": id,
		})
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"address_id": id,
	})
	return nil
}
