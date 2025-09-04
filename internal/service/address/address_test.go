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
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressService_Create(t *testing.T) {
	t.Run("sucesso na criação do endereço", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(1)
		addressModel := &model.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Address) bool {
			return m.UserID != nil &&
				*m.UserID == *addressModel.UserID &&
				m.Street == addressModel.Street &&
				m.City == addressModel.City &&
				m.State == addressModel.State &&
				m.Country == addressModel.Country &&
				m.PostalCode == addressModel.PostalCode
		})).Return(addressModel, nil)

		createdAddress, err := service.Create(context.Background(), addressModel)

		assert.NoError(t, err)
		assert.Equal(t, addressModel.Street, createdAddress.Street)
		assert.Equal(t, addressModel.City, createdAddress.City)
		assert.Equal(t, addressModel.State, createdAddress.State)
		assert.Equal(t, addressModel.Country, createdAddress.Country)
		assert.Equal(t, addressModel.PostalCode, createdAddress.PostalCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado ao buscar endereço", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressID := int64(1)
		expectedErr := err_msg.ErrGet

		mockRepo.
			On("GetByID", mock.Anything, addressID).
			Return((*model.Address)(nil), expectedErr)

		result, err := service.GetByID(context.Background(), addressID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), err_msg.ErrGet.Error())
		assert.Contains(t, err.Error(), expectedErr.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha na validação do endereço UserID/ClientID/SupplierID obrigatório", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		// sem IDs
		addressModel := &model.Address{
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		createdAddress, err := service.Create(context.Background(), addressModel)

		assert.Nil(t, createdAddress)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("falha no repositório ao criar endereço", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(1)
		addressModel := &model.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		expectedErr := errors.New("erro no banco")

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Address) bool {
			return m.UserID != nil &&
				*m.UserID == *addressModel.UserID &&
				m.Street == addressModel.Street &&
				m.City == addressModel.City &&
				m.State == addressModel.State &&
				m.Country == addressModel.Country &&
				m.PostalCode == addressModel.PostalCode
		})).Return((*model.Address)(nil), expectedErr)

		createdAddress, err := service.Create(context.Background(), addressModel)

		assert.Nil(t, createdAddress)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByID(t *testing.T) {
	mockRepo := new(mock_address.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereço por ID", func(t *testing.T) {
		addressModel := &model.Address{
			ID:         1,
			UserID:     nil,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "Estado Teste",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(addressModel, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, addressModel.ID, result.ID)
		assert.Equal(t, addressModel.Street, result.Street)
		assert.Equal(t, addressModel.City, result.City)
		assert.Equal(t, addressModel.State, result.State)
		assert.Equal(t, addressModel.Country, result.Country)
		assert.Equal(t, addressModel.PostalCode, result.PostalCode)

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar endereço com ID inválido", func(t *testing.T) {
		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrID.Error())
	})

	t.Run("endereço não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return((*model.Address)(nil), err_msg.ErrNotFound).Once()

		result, err := service.GetByID(context.Background(), 2)

		assert.ErrorIs(t, err, err_msg.ErrNotFound)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByUserID(t *testing.T) {
	mockRepo := new(mock_address.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereços por UserID", func(t *testing.T) {
		addressModels := []*models.Address{
			{
				ID:         1,
				UserID:     utils.Int64Ptr(1),
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "Estado Teste",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
		}

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(addressModels, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Rua Teste", result[0].Street)
		assert.Equal(t, "Cidade Teste", result[0].City)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereços com UserID inválido", func(t *testing.T) {
		service := NewAddressService(nil)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrID.Error())
	})

	t.Run("nenhum endereço encontrado por UserID", func(t *testing.T) {
		mockRepo.On("GetByUserID", mock.Anything, int64(2)).Return(nil, err_msg.ErrNotFound)

		result, err := service.GetByUserID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, err_msg.ErrNotFound.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByClientID(t *testing.T) {
	mockRepo := new(mock_address.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereços por ClientID", func(t *testing.T) {
		addressesModel := []*models.Address{
			{
				ID:         1,
				ClientID:   utils.Int64Ptr(1),
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "Estado Teste",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
		}

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).Return(addressesModel, nil)

		result, err := service.GetByClientID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, addressesModel, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar endereços com ClientID inválido", func(t *testing.T) {
		result, err := service.GetByClientID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrID.Error())
	})

	t.Run("nenhum endereço encontrado por ClientID", func(t *testing.T) {
		mockRepo.On("GetByClientID", mock.Anything, int64(2)).Return(nil, err_msg.ErrNotFound)

		result, err := service.GetByClientID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, err_msg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetBySupplierID(t *testing.T) {
	mockRepo := new(mock_address.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereços por SupplierID", func(t *testing.T) {
		addressesModel := []*models.Address{
			{
				ID:         1,
				SupplierID: utils.Int64Ptr(1),
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "Estado Teste",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return(addressesModel, nil)

		result, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, addressesModel, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar endereços com SupplierID inválido", func(t *testing.T) {
		result, err := service.GetBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrID.Error())
	})

	t.Run("nenhum endereço encontrado por SupplierID", func(t *testing.T) {
		mockRepo.On("GetBySupplierID", mock.Anything, int64(2)).Return(nil, err_msg.ErrNotFound)

		result, err := service.GetBySupplierID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, err_msg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_UpdateAddress(t *testing.T) {
	makeAddressModel := func() *models.Address {
		return &models.Address{
			ID:         1,
			UserID:     utils.Int64Ptr(1),
			Street:     "Nova Rua",
			City:       "Nova Cidade",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}
	}

	t.Run("sucesso na atualização do endereço", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := makeAddressModel()

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *models.Address) bool {
			return a != nil && a.ID == addressModel.ID && *a.UserID == *addressModel.UserID
		})).Return(nil)

		err := service.Update(context.Background(), addressModel)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha na validação do endereço no update", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := &models.Address{
			ID:     1,
			Street: "", // inválido porque não tem UserID/ClientID/SupplierID
		}

		err := service.Update(context.Background(), addressModel)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("erro genérico ao atualizar endereço", func(t *testing.T) {
		mockRepo := new(mock_address.MockAddressRepository)
		service := NewAddressService(mockRepo)

		addressModel := &models.Address{
			ID:         1,
			UserID:     utils.Int64Ptr(1),
			Street:     "Rua Erro Genérico",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "00456000",
		}

		mockRepo.On("Update", mock.Anything, addressModel).Return(fmt.Errorf("erro inesperado no banco"))

		err := service.Update(context.Background(), addressModel)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar")
		assert.ErrorContains(t, err, "erro inesperado no banco")
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
