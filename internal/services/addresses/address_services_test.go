package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper para ponteiro
func int64Ptr(i int64) *int64 {
	return &i
}

func TestAddressService_Create(t *testing.T) {
	t.Run("sucesso na criação do endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New()) // logger real ou mockado se preferir

		service := NewAddressService(mockRepo, logger)

		userID := int64(1)
		address := &models.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		mockRepo.On("Create", mock.Anything, address).Return(address, nil)

		createdAddress, err := service.Create(context.Background(), address)

		assert.NoError(t, err)
		assert.Equal(t, address, createdAddress)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado ao buscar endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New())
		service := NewAddressService(mockRepo, logger)

		addressID := int64(1)
		expectedErr := errors.New("erro de banco de dados")

		mockRepo.
			On("GetByID", mock.Anything, addressID).
			Return((*models.Address)(nil), expectedErr)

		result, err := service.GetByID(context.Background(), addressID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrCheckAddress.Error())
		assert.Contains(t, err.Error(), expectedErr.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha na validação do endereço - UserID/ClientID/SupplierID obrigatório", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New())

		service := NewAddressService(mockRepo, logger)

		address := &models.Address{
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		createdAddress, err := service.Create(context.Background(), address)

		assert.Nil(t, createdAddress)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("falha no repositório ao criar endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New())

		service := NewAddressService(mockRepo, logger)

		userID := int64(1)
		address := &models.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Create", mock.Anything, address).Return((*models.Address)(nil), expectedErr)

		createdAddress, err := service.Create(context.Background(), address)

		assert.Nil(t, createdAddress)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByID(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	mockLogger := logger.NewLoggerAdapter(logrus.New())
	service := NewAddressService(mockRepo, mockLogger)

	t.Run("sucesso ao buscar endereço por ID", func(t *testing.T) {
		address := &models.Address{
			ID:         1,
			UserID:     nil,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "Estado Teste",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(address, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, address, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar endereço com ID inválido", func(t *testing.T) {
		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("endereço não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return((*models.Address)(nil), repositories.ErrAddressNotFound).Once() // Correção aqui

		result, err := service.GetByID(context.Background(), 2)

		assert.ErrorIs(t, err, ErrAddressNotFound)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

}

func TestAddressService_GetByUserID(t *testing.T) {
	mockLogger := logger.NewLoggerAdapter(logrus.New())
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo, mockLogger)

	t.Run("sucesso ao buscar endereços por UserID", func(t *testing.T) {
		addresses := []*models.Address{
			{
				ID:         1,
				UserID:     int64Ptr(1),
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "Estado Teste",
				Country:    "Brasil",
				PostalCode: "12345-678",
			},
		}

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(addresses, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, addresses, result)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereços com UserID inválido", func(t *testing.T) {
		service := NewAddressService(nil, mockLogger)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("nenhum endereço encontrado por UserID", func(t *testing.T) {
		mockRepo.On("GetByUserID", mock.Anything, int64(2)).Return(nil, ErrAddressNotFound)

		result, err := service.GetByUserID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, ErrAddressNotFound.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByClientID(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	logger := logger.NewLoggerAdapter(logrus.New()) // ou use um logger mock, se preferir
	service := NewAddressService(mockRepo, logger)

	t.Run("sucesso ao buscar endereços por ClientID", func(t *testing.T) {
		clientID := int64(2)
		addresses := []*models.Address{
			{
				ID:         1,
				ClientID:   &clientID,
				Street:     "Rua Cliente",
				City:       "Cidade Cliente",
				State:      "Estado Cliente",
				Country:    "Brasil",
				PostalCode: "98765-432",
			},
		}

		mockRepo.On("GetByClientID", mock.Anything, clientID).Return(addresses, nil)

		result, err := service.GetByClientID(context.Background(), clientID)

		assert.NoError(t, err)
		assert.Equal(t, addresses, result)
		mockRepo.AssertExpectations(t)
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereços com ClientID inválido", func(t *testing.T) {
		service := NewAddressService(mockRepo, logger) // ou nil se quiser simular service sem dependência

		result, err := service.GetByClientID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("nenhum endereço encontrado por ClientID", func(t *testing.T) {
		clientID := int64(2)
		mockRepo.On("GetByClientID", mock.Anything, clientID).Return(nil, ErrAddressNotFound)

		result, err := service.GetByClientID(context.Background(), clientID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, ErrAddressNotFound.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetBySupplierID(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	logger := logger.NewLoggerAdapter(logrus.New())
	service := NewAddressService(mockRepo, logger)

	t.Run("sucesso ao buscar endereços por SupplierID", func(t *testing.T) {
		supplierID := int64(5)
		addresses := []*models.Address{
			{
				ID:         3,
				SupplierID: &supplierID,
				Street:     "Rua Fornecedor",
				City:       "Cidade Fornecedor",
				State:      "Estado Fornecedor",
				Country:    "Brasil",
				PostalCode: "54321-000",
			},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, supplierID).Return(addresses, nil)

		result, err := service.GetBySupplierID(context.Background(), supplierID)

		assert.NoError(t, err)
		assert.Equal(t, addresses, result)
		mockRepo.AssertExpectations(t)
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereços com SupplierID inválido", func(t *testing.T) {
		service := NewAddressService(mockRepo, logger)

		result, err := service.GetBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrAddressIDRequired)
	})

	t.Run("nenhum endereço encontrado por SupplierID", func(t *testing.T) {
		supplierID := int64(5)
		mockRepo.On("GetBySupplierID", mock.Anything, supplierID).Return(nil, ErrAddressNotFound)

		result, err := service.GetBySupplierID(context.Background(), supplierID)

		assert.ErrorIs(t, err, ErrAddressNotFound)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_UpdateAddress(t *testing.T) {
	makeAddress := func() models.Address {
		userID := int64(1)
		return models.Address{
			ID:         1,
			UserID:     &userID,
			Street:     "Nova Rua",
			City:       "Nova Cidade",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "99999-999",
		}
	}

	t.Run("sucesso na atualização do endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New())
		service := NewAddressService(mockRepo, logger)

		address := makeAddress()

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *models.Address) bool {
			return a != nil && a.ID == address.ID && *a.UserID == *address.UserID
		})).Return(nil)

		err := service.Update(context.Background(), &address)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao atualizar endereço com ID inválido", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New())
		service := NewAddressService(mockRepo, logger)

		userID := int64(1)
		address := models.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		err := service.Update(context.Background(), &address)
		assert.ErrorIs(t, err, ErrAddressIDRequired)
	})

	t.Run("falha na validação do endereço no update", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New())
		service := NewAddressService(mockRepo, logger)

		address := &models.Address{
			ID:     1,
			Street: "", // Campos inválidos
		}

		err := service.Update(context.Background(), address)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("erro genérico ao atualizar endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		logger := logger.NewLoggerAdapter(logrus.New())
		service := NewAddressService(mockRepo, logger)

		userID := int64(1)
		address := &models.Address{
			ID:         1,
			UserID:     &userID,
			Street:     "Rua Erro Genérico",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "00456-000",
		}

		mockRepo.On("Update", mock.Anything, address).Return(fmt.Errorf("erro inesperado no banco"))

		err := service.Update(context.Background(), address)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "address: erro ao atualizar")
		assert.ErrorContains(t, err, "erro inesperado no banco")
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_DeleteAddress(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	logger := logger.NewLoggerAdapter(logrus.New())
	service := NewAddressService(mockRepo, logger)

	t.Run("sucesso ao deletar endereço", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("erro ao deletar com ID inválido", func(t *testing.T) {
		err := service.Delete(context.Background(), 0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrAddressIDRequired)
	})

	t.Run("erro ao deletar do repositório", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(fmt.Errorf("erro no banco"))

		err := service.Delete(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "address: erro ao deletar")
		assert.ErrorContains(t, err, "erro no banco")
		mockRepo.AssertExpectations(t)
	})
}
