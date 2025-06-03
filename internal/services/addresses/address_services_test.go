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

var ErrAddressNotFound = errors.New("address: endereço não encontrado")

func TestAddressService_Create(t *testing.T) {
	mockRepo := new(repositories.MockAddressRepository)
	service := NewAddressService(mockRepo)

	t.Run("sucesso na criação do endereço", func(t *testing.T) {
		address := &models.Address{
			ID:         0,
			UserID:     nil,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "Estado Teste",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		mockRepo.On("Create", mock.Anything, address).Return(address, nil)

		createdAddress, err := service.Create(context.Background(), address)

		assert.NoError(t, err)
		assert.Equal(t, address, createdAddress)
		mockRepo.AssertExpectations(t)
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

func TestAddressService_UpdateAddress(t *testing.T) {
	makeAddress := func() models.Address {
		return models.Address{
			ID:         1,
			Street:     "Nova Rua",
			City:       "Nova Cidade",
			State:      "Novo Estado",
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
			return a != nil && a.ID != 0 && a.ID == address.ID && a.Street == address.Street && a.Version == address.Version
		})).Return(nil)

		err := service.Update(context.Background(), &address)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao atualizar endereço com ID inválido", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		address := models.Address{
			Street:  "Rua Teste",
			Version: 1,
		}

		err := service.Update(context.Background(), &address)

		assert.ErrorIs(t, err, ErrAddressIDRequired)
	})

	t.Run("erro ao atualizar endereço com versão zero", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		address := models.Address{
			ID:      1,
			Street:  "Rua Teste",
			Version: 0,
		}

		err := service.Update(context.Background(), &address)

		assert.ErrorContains(t, err, "versão obrigatória")
	})

	t.Run("erro por conflito de versão", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		address := &models.Address{
			ID:         1,
			Street:     "Rua Conflito",
			City:       "Cidade Conflito",
			State:      "Estado Conflito",
			Country:    "Brasil",
			PostalCode: "00000-000",
			Version:    2,
		}

		mockRepo.On("Update", mock.Anything, mock.Anything).
			Return(repositories.ErrVersionConflict)

		err := service.Update(context.Background(), address)

		assert.ErrorIs(t, err, ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro genérico ao atualizar endereço", func(t *testing.T) {
		mockRepo := new(repositories.MockAddressRepository)
		service := NewAddressService(mockRepo)

		address := &models.Address{
			ID:         1,
			Street:     "Rua Erro Genérico",
			City:       "Cidade Teste",
			State:      "Estado Teste",
			Country:    "Brasil",
			PostalCode: "00000-000",
			Version:    1,
		}

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *models.Address) bool {
			return a.ID != 0 && a.ID == address.ID
		})).Return(fmt.Errorf("erro inesperado no banco"))

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
		mockRepo.On("Delete", mock.Anything, int64(1), 2).Return(nil)

		err := service.Delete(context.Background(), int64(1), 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao deletar com ID inválido", func(t *testing.T) {
		err := service.Delete(context.Background(), 0, 2)

		assert.Error(t, err)
		assert.EqualError(t, err, ErrAddressIDRequired.Error())
	})

	t.Run("erro ao deletar com versão inválida", func(t *testing.T) {
		err := service.Delete(context.Background(), 1, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, ErrVersionRequired.Error())
	})
}
