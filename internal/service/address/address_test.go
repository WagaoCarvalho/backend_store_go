package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mock_address "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/address"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressService_Create(t *testing.T) {
	userID := int64(1)

	t.Run("falha na validacao do endereco", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		// endereço inválido: nenhum ID preenchido, campos obrigatórios vazios
		addressModel := &model.Address{}

		createdAddress, err := service.Create(context.Background(), addressModel)

		assert.Nil(t, createdAddress)
		assert.Error(t, err)
		// ajustado para matchar o que o model realmente retorna
		assert.ErrorContains(t, err, "user_id/client_id/supplier_id")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("sucesso na criação do endereço", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := &model.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
			IsActive:   true,
		}

		mockRepo.On("Create", mock.Anything, addressModel).Return(addressModel, nil)

		createdAddress, err := service.Create(context.Background(), addressModel)

		assert.NoError(t, err)
		assert.Equal(t, addressModel, createdAddress)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório ao criar endereço", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := &model.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
			IsActive:   true,
		}

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Create", mock.Anything, addressModel).Return((*model.Address)(nil), expectedErr)

		createdAddress, err := service.Create(context.Background(), addressModel)

		assert.Nil(t, createdAddress)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrID)
		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("endereço não encontrado", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressID := int64(1)
		mockRepo.On("GetByID", mock.Anything, addressID).Return((*models.Address)(nil), err_msg.ErrNotFound)

		result, err := service.GetByID(context.Background(), addressID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado do repositório", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressID := int64(2)
		unexpectedErr := errors.New("erro no banco")
		mockRepo.On("GetByID", mock.Anything, addressID).Return((*models.Address)(nil), unexpectedErr)

		result, err := service.GetByID(context.Background(), addressID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), err_msg.ErrGet.Error())
		assert.Contains(t, err.Error(), unexpectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso na busca", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressID := int64(3)
		expectedAddress := &models.Address{
			ID:         addressID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		mockRepo.On("GetByID", mock.Anything, addressID).Return(expectedAddress, nil)

		result, err := service.GetByID(context.Background(), addressID)

		assert.NoError(t, err)
		assert.Equal(t, expectedAddress, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByUserID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrID)
		mockRepo.AssertNotCalled(t, "GetByUserID", mock.Anything, mock.Anything)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(1)
		expectedErr := errors.New("erro no banco")

		mockRepo.On("GetByUserID", mock.Anything, userID).Return(nil, expectedErr)

		result, err := service.GetByUserID(context.Background(), userID)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso na busca", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(2)
		expectedAddresses := []*models.Address{
			{ID: 1, UserID: &userID, Street: "Rua A", City: "Cidade A", State: "SP", Country: "Brasil", PostalCode: "12345678"},
			{ID: 2, UserID: &userID, Street: "Rua B", City: "Cidade B", State: "RJ", Country: "Brasil", PostalCode: "87654321"},
		}

		mockRepo.On("GetByUserID", mock.Anything, userID).Return(expectedAddresses, nil)

		result, err := service.GetByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedAddresses, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByClientID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		result, err := service.GetByClientID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrID)
		mockRepo.AssertNotCalled(t, "GetByClientID", mock.Anything, mock.Anything)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		clientID := int64(1)
		expectedErr := errors.New("erro no banco")

		mockRepo.On("GetByClientID", mock.Anything, clientID).Return(nil, expectedErr)

		result, err := service.GetByClientID(context.Background(), clientID)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso na busca", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		clientID := int64(2)
		expectedAddresses := []*models.Address{
			{ID: 1, ClientID: &clientID, Street: "Rua A", City: "Cidade A", State: "SP", Country: "Brasil", PostalCode: "12345678"},
			{ID: 2, ClientID: &clientID, Street: "Rua B", City: "Cidade B", State: "RJ", Country: "Brasil", PostalCode: "87654321"},
		}

		mockRepo.On("GetByClientID", mock.Anything, clientID).Return(expectedAddresses, nil)

		result, err := service.GetByClientID(context.Background(), clientID)

		assert.NoError(t, err)
		assert.Equal(t, expectedAddresses, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetBySupplierID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		result, err := service.GetBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrID)
		mockRepo.AssertNotCalled(t, "GetBySupplierID", mock.Anything, mock.Anything)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		supplierID := int64(1)
		expectedErr := errors.New("erro no banco")

		mockRepo.On("GetBySupplierID", mock.Anything, supplierID).Return(nil, expectedErr)

		result, err := service.GetBySupplierID(context.Background(), supplierID)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso na busca", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		supplierID := int64(2)
		expectedAddresses := []*models.Address{
			{ID: 1, SupplierID: &supplierID, Street: "Rua A", City: "Cidade A", State: "SP", Country: "Brasil", PostalCode: "12345678"},
			{ID: 2, SupplierID: &supplierID, Street: "Rua B", City: "Cidade B", State: "RJ", Country: "Brasil", PostalCode: "87654321"},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, supplierID).Return(expectedAddresses, nil)

		result, err := service.GetBySupplierID(context.Background(), supplierID)

		assert.NoError(t, err)
		assert.Equal(t, expectedAddresses, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_Update(t *testing.T) {
	userID := int64(1)

	t.Run("falha na validacao do endereco", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := &model.Address{} // endereço inválido, sem IDs nem campos obrigatórios

		err := service.Update(context.Background(), addressModel)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "user_id/client_id/supplier_id")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("sucesso no update do endereco", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := &model.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
			IsActive:   true,
		}

		mockRepo.On("Update", mock.Anything, addressModel).Return(nil)

		err := service.Update(context.Background(), addressModel)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositorio ao atualizar endereco", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := &model.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
			IsActive:   true,
		}

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Update", mock.Anything, addressModel).Return(expectedErr)

		err := service.Update(context.Background(), addressModel)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), err_msg.ErrUpdate.Error())
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_DeleteAddress(t *testing.T) {
	mockRepo := new(mock_address.MockAddressRepository)

	service := NewAddressService(mockRepo)

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
		assert.ErrorIs(t, err, err_msg.ErrID)
	})

	t.Run("erro ao deletar do repositório", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(fmt.Errorf("erro no banco"))

		err := service.Delete(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao deletar")
		assert.ErrorContains(t, err, "erro no banco")
		mockRepo.AssertExpectations(t)
	})
}
