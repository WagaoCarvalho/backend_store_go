package services

import (
	"context"
	"errors"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client_cpf"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_GetByID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("falha ao buscar cliente", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("GetByID", mock.Anything, int64(1)).
			Return(nil, expectedErr)

		result, err := service.GetByID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := &models.ClientCpf{ID: 1, Name: "teste"}

		mockRepo.
			On("GetByID", mock.Anything, int64(1)).
			Return(client, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, client, result)
		mockRepo.AssertExpectations(t)
	})
}
