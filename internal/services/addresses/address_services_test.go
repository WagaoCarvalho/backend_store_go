package services_test

import (
	"context"
	"errors"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAddressRepository struct {
	mock.Mock
}

func (m *MockAddressRepository) Create(ctx context.Context, address models.Address) (models.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetByID(ctx context.Context, id int) (models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressRepository) Update(ctx context.Context, address models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAddressService_Create(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	t.Run("sucesso na criação do endereço", func(t *testing.T) {
		address := models.Address{
			ID:         nil,
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

	t.Run("falha na criação do endereço com dados inválidos", func(t *testing.T) {
		address := models.Address{
			Street: "",
			City:   "Cidade Teste",
			State:  "Estado Teste",
		}

		createdAddress, err := service.Create(context.Background(), address)

		assert.Error(t, err)
		assert.Equal(t, "dados do endereço inválidos", err.Error())
		assert.Equal(t, models.Address{}, createdAddress)

		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestAddressService_GetByID(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	t.Run("sucesso ao buscar endereço por ID", func(t *testing.T) {
		address := models.Address{
			ID:         nil,
			UserID:     nil,
			Street:     "Rua Teste",
			City:       "Cidade Teste",
			State:      "Estado Teste",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		mockRepo.On("GetByID", mock.Anything, 1).Return(address, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, address, result)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil // limpa mocks para o próximo subteste
		mockRepo.Calls = nil
	})

	t.Run("endereço não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, 1).Return(models.Address{}, errors.New("endereço não encontrado"))

		result, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, models.Address{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_UpdateAddress(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	t.Run("sucesso na atualização do endereço", func(t *testing.T) {
		id := int64(1)
		address := models.Address{
			ID:         &id,
			Street:     "Nova Rua",
			City:       "Nova Cidade",
			State:      "Novo Estado",
			Country:    "Brasil",
			PostalCode: "99999-999",
		}

		mockRepo.On("Update", mock.Anything, address).Return(nil)

		err := service.Update(context.Background(), address)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao atualizar endereço com ID inválido", func(t *testing.T) {
		address := models.Address{
			Street: "Nova Rua",
		}

		err := service.Update(context.Background(), address)

		assert.Error(t, err)
		assert.Equal(t, "ID do endereço é obrigatório", err.Error())
	})
}

func TestAddressService_DeleteAddress(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	t.Run("sucesso ao deletar endereço", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, 1).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao deletar com ID inválido", func(t *testing.T) {
		err := service.Delete(context.Background(), 0)

		assert.Error(t, err)
		assert.Equal(t, "ID do endereço é obrigatório", err.Error())
	})
}
