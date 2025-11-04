package services

import (
	"context"
	"errors"
	"testing"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressService_DisableAddress(t *testing.T) {
	mockRepoAddress := new(mockAddress.MockAddress)
	mockRepoClient := new(mockClient.MockClient)
	mockRepoUser := new(mockUser.MockUser)
	mockSupplier := new(mockSupplier.MockSupplier)
	service := NewAddress(mockRepoAddress, mockRepoClient, mockRepoUser, mockSupplier)

	t.Run("sucesso ao desabilitar endereço", func(t *testing.T) {
		mockRepoAddress.On("Disable", mock.Anything, int64(1)).Return(nil)

		err := service.Disable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepoAddress.AssertExpectations(t)

		mockRepoAddress.ExpectedCalls = nil
		mockRepoAddress.Calls = nil
	})

	t.Run("erro ao desabilitar com ID inválido", func(t *testing.T) {
		err := service.Disable(context.Background(), 0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro ao desabilitar no repository", func(t *testing.T) {
		mockRepoAddress.On("Disable", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.Disable(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDisable)
		mockRepoAddress.AssertExpectations(t)

		mockRepoAddress.ExpectedCalls = nil
		mockRepoAddress.Calls = nil
	})
}

func TestAddressService_EnableAddress(t *testing.T) {
	mockRepoAddress := new(mockAddress.MockAddress)
	mockRepoClient := new(mockClient.MockClient)
	mockRepoUser := new(mockUser.MockUser)
	mockSupplier := new(mockSupplier.MockSupplier)
	service := NewAddress(mockRepoAddress, mockRepoClient, mockRepoUser, mockSupplier)

	t.Run("sucesso ao habilitar endereço", func(t *testing.T) {
		mockRepoAddress.On("Enable", mock.Anything, int64(1)).Return(nil)

		err := service.Enable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepoAddress.AssertExpectations(t)

		mockRepoAddress.ExpectedCalls = nil
		mockRepoAddress.Calls = nil
	})

	t.Run("erro ao habilitar com ID inválido", func(t *testing.T) {
		err := service.Enable(context.Background(), 0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro ao habilitar no repository", func(t *testing.T) {
		mockRepoAddress.On("Enable", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.Enable(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrEnable)
		mockRepoAddress.AssertExpectations(t)

		mockRepoAddress.ExpectedCalls = nil
		mockRepoAddress.Calls = nil
	})
}
