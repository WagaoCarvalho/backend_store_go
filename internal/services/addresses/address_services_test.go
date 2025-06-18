package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper para ponteiro
func int64Ptr(i int64) *int64 {
	return &i
}

func TestAddressService_Create(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo)
	t.Run("sucesso na criação do endereço", func(t *testing.T) {
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

	t.Run("falha na validação do endereço - UserID/ClientID/SupplierID obrigatório", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		address := &models.Address{
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "Estado Teste",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		createdAddress, err := service.Create(context.Background(), address)

		assert.Nil(t, createdAddress)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

}

func TestAddressService_GetByID(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereço por ID", func(t *testing.T) {
		address := &models.Address{
			ID:         0,
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

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereço com ID inválido", func(t *testing.T) {
		service := NewAddressService(nil)

		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("endereço não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&models.Address{}, ErrAddressNotFound)

		result, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, &models.Address{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByUserID(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereço por UserID", func(t *testing.T) {
		address := &models.Address{
			ID:         1,
			UserID:     int64Ptr(1),
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "Estado Teste",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(address, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, address, result)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereço com UserID inválido", func(t *testing.T) {
		service := NewAddressService(nil)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("endereço não encontrado por UserID", func(t *testing.T) {
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(&models.Address{}, ErrAddressNotFound)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, &models.Address{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByClientID(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereço por ClientID", func(t *testing.T) {
		address := &models.Address{
			ID:         1,
			ClientID:   int64Ptr(2),
			Street:     "Rua Cliente",
			City:       "Cidade Cliente",
			State:      "Estado Cliente",
			Country:    "Brasil",
			PostalCode: "98765-432",
		}

		mockRepo.On("GetByClientID", mock.Anything, int64(2)).Return(address, nil)

		result, err := service.GetByClientID(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, address, result)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereço com ClientID inválido", func(t *testing.T) {
		service := NewAddressService(nil)

		result, err := service.GetByClientID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("endereço não encontrado por ClientID", func(t *testing.T) {
		mockRepo.On("GetByClientID", mock.Anything, int64(2)).Return(&models.Address{}, ErrAddressNotFound)

		result, err := service.GetByClientID(context.Background(), 2)

		assert.Error(t, err)
		assert.Equal(t, &models.Address{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetBySupplierID(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereço por SupplierID", func(t *testing.T) {
		address := &models.Address{
			ID:         3,
			SupplierID: int64Ptr(5),
			Street:     "Rua Fornecedor",
			City:       "Cidade Fornecedor",
			State:      "Estado Fornecedor",
			Country:    "Brasil",
			PostalCode: "54321-000",
		}

		mockRepo.On("GetBySupplierID", mock.Anything, int64(5)).Return(address, nil)

		result, err := service.GetBySupplierID(context.Background(), 5)

		assert.NoError(t, err)
		assert.Equal(t, address, result)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar endereço com SupplierID inválido", func(t *testing.T) {
		service := NewAddressService(nil)

		result, err := service.GetBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("endereço não encontrado por SupplierID", func(t *testing.T) {
		mockRepo.On("GetBySupplierID", mock.Anything, int64(5)).Return(&models.Address{}, ErrAddressNotFound)

		result, err := service.GetBySupplierID(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, &models.Address{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetVersionByID(t *testing.T) {
	t.Run("sucesso ao buscar versão do endereço por ID", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		expectedVersion := 3
		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(expectedVersion, nil)

		result, err := service.GetVersionByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar versão com ID inválido", func(t *testing.T) {
		service := NewAddressService(nil)

		result, err := service.GetVersionByID(context.Background(), 0)

		assert.Equal(t, 0, result)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidID.Error())
	})

	t.Run("endereço não encontrado ao buscar versão", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(0, repositories.ErrAddressNotFound)

		result, err := service.GetVersionByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, 0, result)
		assert.EqualError(t, err, ErrAddressNotFound.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar versão no repositório", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		expectedErr := errors.New("erro no banco de dados")
		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(0, expectedErr)

		result, err := service.GetVersionByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, 0, result)
		assert.Contains(t, err.Error(), "failed to get address version")
		assert.True(t, errors.Is(err, expectedErr), "deve envolver o erro original")
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
			Version:    1,
		}
	}

	t.Run("sucesso na atualização do endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		address := makeAddress()

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *models.Address) bool {
			return a != nil && a.ID == address.ID && a.Version == address.Version &&
				a.Street == address.Street && a.UserID != nil && *a.UserID == *address.UserID
		})).Return(nil)

		err := service.Update(context.Background(), &address)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao atualizar endereço com ID inválido", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(1)
		address := models.Address{
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
			Version:    1,
		}

		err := service.Update(context.Background(), &address)

		assert.ErrorIs(t, err, ErrAddressIDRequired)
	})

	t.Run("erro ao atualizar endereço com versão zero", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(1)
		address := models.Address{
			ID:         1,
			UserID:     &userID,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
			Version:    0,
		}

		err := service.Update(context.Background(), &address)

		assert.ErrorIs(t, err, ErrVersionRequired)
	})

	t.Run("falha na validação do endereço no update", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		address := &models.Address{
			ID:      1,
			Version: 1,
			Street:  "",
		}

		err := service.Update(context.Background(), address)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("erro por conflito de versão", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(1)
		address := &models.Address{
			ID:         1,
			UserID:     &userID,
			Street:     "Rua Conflito",
			City:       "Cidade Conflito",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
			Version:    2,
		}

		mockRepo.On("Update", mock.Anything, address).
			Return(repositories.ErrVersionConflict)

		err := service.Update(context.Background(), address)

		assert.ErrorIs(t, err, ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro genérico ao atualizar endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		userID := int64(1)
		address := &models.Address{
			ID:         1,
			UserID:     &userID,
			Street:     "Rua Erro Genérico",
			City:       "Cidade Teste",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "00456-000",
			Version:    1,
		}

		mockRepo.On("Update", mock.Anything, address).
			Return(fmt.Errorf("erro inesperado no banco"))

		err := service.Update(context.Background(), address)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "address: erro ao atualizar")
		assert.ErrorContains(t, err, "erro inesperado no banco")
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_DeleteAddress(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso ao deletar endereço", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), int64(1))

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao deletar com ID inválido", func(t *testing.T) {
		err := service.Delete(context.Background(), 0)

		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})
}
