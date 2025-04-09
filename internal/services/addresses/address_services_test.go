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

func (m *MockAddressRepository) CreateAddress(ctx context.Context, address models.Address) (models.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressRepository) GetAddressByID(ctx context.Context, id int) (models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressRepository) UpdateAddress(ctx context.Context, address models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressRepository) DeleteAddress(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAddressService_CreateAddress_Success(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	address := models.Address{
		ID:         1,
		UserID:     nil,
		Street:     "Rua Teste",
		City:       "Cidade Teste",
		State:      "Estado Teste",
		Country:    "Brasil",
		PostalCode: "12345-678",
	}

	mockRepo.On("CreateAddress", mock.Anything, address).Return(address, nil)

	createdAddress, err := service.CreateAddress(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, address, createdAddress)
	mockRepo.AssertExpectations(t)
}

func TestAddressService_CreateAddress_Error(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	address := models.Address{
		Street: "",
		City:   "Cidade Teste",
		State:  "Estado Teste",
	}

	createdAddress, err := service.CreateAddress(context.Background(), address)

	assert.Error(t, err)
	assert.Equal(t, "dados do endereço inválidos", err.Error())
	assert.Equal(t, models.Address{}, createdAddress)

	mockRepo.AssertNotCalled(t, "CreateAddress")
}

func TestAddressService_GetAddressByID_Success(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	address := models.Address{
		ID:         1,
		UserID:     nil,
		Street:     "Rua Teste",
		City:       "Cidade Teste",
		State:      "Estado Teste",
		Country:    "Brasil",
		PostalCode: "12345-678",
	}

	mockRepo.On("GetAddressByID", mock.Anything, 1).Return(address, nil)

	result, err := service.GetAddressByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, address, result)
	mockRepo.AssertExpectations(t)
}

func TestAddressService_GetAddressByID_NotFound(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	mockRepo.On("GetAddressByID", mock.Anything, 1).Return(models.Address{}, errors.New("endereço não encontrado"))

	result, err := service.GetAddressByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, models.Address{}, result)
	mockRepo.AssertExpectations(t)
}

func TestAddressService_UpdateAddress_Success(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	address := models.Address{
		ID:         1,
		Street:     "Nova Rua",
		City:       "Nova Cidade",
		State:      "Novo Estado",
		Country:    "Brasil",
		PostalCode: "99999-999",
	}

	mockRepo.On("UpdateAddress", mock.Anything, address).Return(nil)

	err := service.UpdateAddress(context.Background(), address)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddressService_UpdateAddress_InvalidID(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	address := models.Address{
		Street: "Nova Rua",
	}

	err := service.UpdateAddress(context.Background(), address)

	assert.Error(t, err)
	assert.Equal(t, "ID do endereço é obrigatório", err.Error())
}

func TestAddressService_DeleteAddress_Success(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	mockRepo.On("DeleteAddress", mock.Anything, 1).Return(nil)

	err := service.DeleteAddress(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddressService_DeleteAddress_InvalidID(t *testing.T) {
	mockRepo := new(MockAddressRepository)
	service := services.NewAddressService(mockRepo)

	err := service.DeleteAddress(context.Background(), 0)

	assert.Error(t, err)
	assert.Equal(t, "ID do endereço é obrigatório", err.Error())
}
