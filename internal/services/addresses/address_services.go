package services

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
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
	repo   repositories.AddressRepository
	logger *logger.LoggerAdapter
}

func NewAddressService(repo repositories.AddressRepository, logger *logger.LoggerAdapter) AddressService {
	return &addressService{
		repo:   repo,
		logger: logger,
	}
}

func (s *addressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	s.logger.Info(ctx, "[addressService] - "+logger.LogCreateInit, map[string]interface{}{
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
	})

	if err := address.Validate(); err != nil {
		s.logger.Warn(ctx, "[addressService] - "+logger.LogErrorValidate, map[string]interface{}{
			"erro": err.Error(),
		})
		return nil, err
	}

	createdAddress, err := s.repo.Create(ctx, address)
	if err != nil {
		s.logger.Error(ctx, err, "[addressService] - "+logger.LogCreateError, map[string]interface{}{
			"street": address.Street,
		})
		return nil, err
	}

	s.logger.Info(ctx, "[addressService] - "+logger.LogCreateSuccess, map[string]interface{}{
		"address_id": createdAddress.ID,
	})

	return createdAddress, nil
}

func (s *addressService) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	s.logger.Info(ctx, "[addressService] - Iniciando busca de endereço por ID", map[string]interface{}{
		"address_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, "[addressService] - ID do endereço não fornecido", map[string]interface{}{
			"address_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	address, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, "[addressService] - Erro ao buscar endereço por ID", map[string]interface{}{
			"address_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, "[addressService] - Endereço recuperado com sucesso", map[string]interface{}{
		"address_id": address.ID,
	})

	return address, nil
}

func (s *addressService) GetByUserID(ctx context.Context, id int64) ([]*models.Address, error) {
	s.logger.Info(ctx, "[addressService] - Iniciando busca de endereços por UserID", map[string]interface{}{
		"user_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, "[addressService] - ID de usuário inválido para busca de endereços", map[string]interface{}{
			"user_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	addresses, err := s.repo.GetByUserID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, "[addressService] - Erro ao buscar endereços no repositório", map[string]interface{}{
			"user_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, "[addressService] - Endereços recuperados com sucesso", map[string]interface{}{
		"user_id":         id,
		"total_addresses": len(addresses),
	})

	return addresses, nil
}

func (s *addressService) GetByClientID(ctx context.Context, id int64) ([]*models.Address, error) {
	s.logger.Info(ctx, "[addressService] - Iniciando busca de endereços por ClientID", map[string]interface{}{
		"client_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, "[addressService] - ClientID não fornecido", map[string]interface{}{
			"client_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	addresses, err := s.repo.GetByClientID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, "[addressService] - Erro ao buscar endereços por ClientID", map[string]interface{}{
			"client_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, "[addressService] - Endereços recuperados com sucesso", map[string]interface{}{
		"client_id":   id,
		"total_items": len(addresses),
	})

	return addresses, nil
}

func (s *addressService) GetBySupplierID(ctx context.Context, id int64) ([]*models.Address, error) {
	s.logger.Info(ctx, "[addressService] - Iniciando busca de endereços por SupplierID", map[string]interface{}{
		"supplier_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, "[addressService] - SupplierID não fornecido", map[string]interface{}{
			"supplier_id": id,
		})
		return nil, ErrAddressIDRequired
	}

	addresses, err := s.repo.GetBySupplierID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, "[addressService] - Erro ao buscar endereços por SupplierID", map[string]interface{}{
			"supplier_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, "[addressService] - Endereços recuperados com sucesso", map[string]interface{}{
		"supplier_id": id,
		"total_items": len(addresses),
	})

	return addresses, nil
}

func (s *addressService) Update(ctx context.Context, address *models.Address) error {
	s.logger.Info(ctx, "[addressService] - Iniciando atualização de endereço", map[string]interface{}{
		"address_id":  address.ID,
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
	})

	if err := address.Validate(); err != nil {
		s.logger.Warn(ctx, "[addressService] - Validação de endereço falhou", map[string]interface{}{
			"erro": err.Error(),
		})
		return err
	}

	if address.ID == 0 {
		s.logger.Warn(ctx, "[addressService] - ID do endereço não fornecido para atualização", map[string]interface{}{
			"address_id": address.ID,
		})
		return ErrAddressIDRequired
	}

	err := s.repo.Update(ctx, address)
	if err != nil {
		s.logger.Error(ctx, err, "[addressService] - Erro ao atualizar endereço no repositório", map[string]interface{}{
			"address_id": address.ID,
		})
		return fmt.Errorf("%w: %v", ErrUpdateAddress, err)
	}

	s.logger.Info(ctx, "[addressService] - Endereço atualizado com sucesso", map[string]interface{}{
		"address_id": address.ID,
	})

	return nil
}

func (s *addressService) Delete(ctx context.Context, id int64) error {
	s.logger.Info(ctx, "[addressService] - Iniciando exclusão de endereço", map[string]interface{}{
		"address_id": id,
	})

	if id == 0 {
		s.logger.Warn(ctx, "[addressService] - ID do endereço não fornecido para exclusão", map[string]interface{}{
			"address_id": id,
		})
		return ErrAddressIDRequired
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, "[addressService] - Erro ao excluir endereço no repositório", map[string]interface{}{
			"address_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAddress, err)
	}

	s.logger.Info(ctx, "[addressService] - Endereço excluído com sucesso", map[string]interface{}{
		"address_id": id,
	})

	return nil
}
