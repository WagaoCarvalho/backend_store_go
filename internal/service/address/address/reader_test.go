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
