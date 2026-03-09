package services

import (
	"context"
	"errors"
	"testing"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	mockClientCpf "github.com/WagaoCarvalho/backend_store_go/infra/mock/client_cpf"
	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressService_GetByID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		service := NewAddressService(
			new(mockAddress.MockAddress),
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro propagado do repositório", func(t *testing.T) {
		mockRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			mockRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		unexpectedErr := errors.New("db error")

		mockRepo.
			On("GetByID", mock.Anything, int64(1)).
			Return((*models.Address)(nil), unexpectedErr)

		result, err := service.GetByID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, unexpectedErr)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			mockRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		expected := &models.Address{ID: 1}

		mockRepo.
			On("GetByID", mock.Anything, int64(1)).
			Return(expected, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestAddressService_GetByUserID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		service := NewAddressService(
			new(mockAddress.MockAddress),
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro do repositório ao buscar endereços", func(t *testing.T) {
		mockAddressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			mockAddressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		unexpectedErr := errors.New("db error")

		mockAddressRepo.
			On("GetByUserID", mock.Anything, int64(1)).
			Return(nil, unexpectedErr)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, unexpectedErr)
	})

	t.Run("nenhum endereço encontrado - retorna lista vazia", func(t *testing.T) {
		mockAddressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			mockAddressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		mockAddressRepo.
			On("GetByUserID", mock.Anything, int64(1)).
			Return([]*models.Address{}, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("endereços encontrados", func(t *testing.T) {
		mockAddressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			mockAddressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		addresses := []*models.Address{
			{ID: 1, UserID: utils.Int64Ptr(1)},
			{ID: 2, UserID: utils.Int64Ptr(1)},
		}

		mockAddressRepo.
			On("GetByUserID", mock.Anything, int64(1)).
			Return(addresses, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, addresses, result)
	})
}

func TestAddressService_GetByClientCpfID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		service := NewAddressService(
			new(mockAddress.MockAddress),
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		result, err := service.GetByClientCpfID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro do repositório de endereço", func(t *testing.T) {
		addressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			addressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		unexpectedErr := errors.New("db error")

		addressRepo.
			On("GetByClientCpfID", mock.Anything, int64(1)).
			Return(nil, unexpectedErr)

		result, err := service.GetByClientCpfID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, unexpectedErr)
	})

	t.Run("nenhum endereço encontrado - retorna lista vazia", func(t *testing.T) {
		addressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			addressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		addressRepo.
			On("GetByClientCpfID", mock.Anything, int64(1)).
			Return([]*models.Address{}, nil)

		result, err := service.GetByClientCpfID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("endereços encontrados", func(t *testing.T) {
		addressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			addressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		addresses := []*models.Address{
			{ID: 1, ClientCpfID: utils.Int64Ptr(1)},
		}

		addressRepo.
			On("GetByClientCpfID", mock.Anything, int64(1)).
			Return(addresses, nil)

		result, err := service.GetByClientCpfID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, addresses, result)
	})
}

func TestAddressService_GetBySupplierID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		service := NewAddressService(
			new(mockAddress.MockAddress),
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		result, err := service.GetBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro do repositório de endereço", func(t *testing.T) {
		addressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			addressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		unexpectedErr := errors.New("db error")

		addressRepo.
			On("GetBySupplierID", mock.Anything, int64(1)).
			Return(nil, unexpectedErr)

		result, err := service.GetBySupplierID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, unexpectedErr)
	})

	t.Run("nenhum endereço encontrado - retorna lista vazia", func(t *testing.T) {
		addressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			addressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		addressRepo.
			On("GetBySupplierID", mock.Anything, int64(1)).
			Return([]*models.Address{}, nil)

		result, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("endereços encontrados", func(t *testing.T) {
		addressRepo := new(mockAddress.MockAddress)

		service := NewAddressService(
			addressRepo,
			new(mockClientCpf.MockClientCpf),
			new(mockUser.MockUser),
			new(mockSupplier.MockSupplier),
		)

		addresses := []*models.Address{
			{ID: 1, SupplierID: utils.Int64Ptr(1)},
		}

		addressRepo.
			On("GetBySupplierID", mock.Anything, int64(1)).
			Return(addresses, nil)

		result, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, addresses, result)
	})
}
