package services

import (
	"context"
	"errors"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client_cpf"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_Disable(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		err := service.Disable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Disable", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("Disable", mock.Anything, int64(1)).
			Return(expectedErr)

		err := service.Disable(context.Background(), 1)

		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		mockRepo.
			On("Disable", mock.Anything, int64(1)).
			Return(nil)

		err := service.Disable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_Enable(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		err := service.Enable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Enable", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("Enable", mock.Anything, int64(1)).
			Return(expectedErr)

		err := service.Enable(context.Background(), 1)

		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		mockRepo.
			On("Enable", mock.Anything, int64(1)).
			Return(nil)

		err := service.Enable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
