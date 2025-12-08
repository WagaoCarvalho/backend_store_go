package services

import (
	"context"
	"errors"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_GetVersionByID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		version, err := service.GetVersionByID(context.Background(), 0)

		assert.Zero(t, version)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetVersionByID", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("GetVersionByID", mock.Anything, int64(1)).
			Return(0, expectedErr)

		version, err := service.GetVersionByID(context.Background(), 1)

		assert.Zero(t, version)
		assert.ErrorContains(t, err, errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		mockRepo.
			On("GetVersionByID", mock.Anything, int64(1)).
			Return(5, nil)

		version, err := service.GetVersionByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 5, version)
		mockRepo.AssertExpectations(t)
	})
}
