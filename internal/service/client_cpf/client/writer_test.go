package services

import (
	"context"
	"errors"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client_cpf"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func validClient() *models.ClientCpf {
	return &models.ClientCpf{
		ID:      1,
		Name:    "Cliente Teste",
		Email:   *utils.StrToPtr("teste@teste.com"),
		CPF:     "12345678901",
		Version: 1,
		Status:  true,
	}
}

func TestClientService_Create(t *testing.T) {
	t.Run("falha quando client é nil", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		result, err := service.Create(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("falha na validação", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		invalid := &models.ClientCpf{Name: ""}

		result, err := service.Create(context.Background(), invalid)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("erro - cliente duplicado", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()

		mockRepo.
			On("Create", mock.Anything, client).
			Return(nil, errMsg.ErrDuplicate).
			Once()

		result, err := service.Create(context.Background(), client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - chave estrangeira inválida", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()

		mockRepo.
			On("Create", mock.Anything, client).
			Return(nil, errMsg.ErrDBInvalidForeignKey).
			Once()

		result, err := service.Create(context.Background(), client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao criar no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()
		expectedErr := errors.New("erro ao criar")

		mockRepo.
			On("Create", mock.Anything, client).
			Return(nil, expectedErr).
			Once()

		result, err := service.Create(context.Background(), client)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, errMsg.ErrCreate.Error())
		assert.ErrorContains(t, err, expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()

		mockRepo.
			On("Create", mock.Anything, client).
			Return(client, nil).
			Once()

		result, err := service.Create(context.Background(), client)

		assert.NoError(t, err)
		assert.Equal(t, client, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_Update(t *testing.T) {
	t.Run("falha - ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()
		client.ID = 0

		err := service.Update(context.Background(), client)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - versão inválida", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()
		client.Version = 0

		err := service.Update(context.Background(), client)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - validação", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()
		client.Email = ""

		err := service.Update(context.Background(), client)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - erro genérico do repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()
		expectedErr := errors.New("db error")

		mockRepo.
			On("Update", mock.Anything, client).
			Return(expectedErr).
			Once()

		err := service.Update(context.Background(), client)

		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		assert.ErrorContains(t, err, expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha - conflito de versão", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()

		mockRepo.
			On("Update", mock.Anything, client).
			Return(errMsg.ErrVersionConflict).
			Once()

		err := service.Update(context.Background(), client)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha - cliente duplicado", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()

		mockRepo.
			On("Update", mock.Anything, client).
			Return(errMsg.ErrDuplicate).
			Once()

		err := service.Update(context.Background(), client)

		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		client := validClient()

		mockRepo.
			On("Update", mock.Anything, client).
			Return(nil).
			Once()

		err := service.Update(context.Background(), client)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_Delete(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		expectedErr := errors.New("db error")

		mockRepo.
			On("Delete", mock.Anything, int64(1)).
			Return(expectedErr).
			Once()

		err := service.Delete(context.Background(), 1)

		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		assert.ErrorContains(t, err, expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfService(mockRepo)

		mockRepo.
			On("Delete", mock.Anything, int64(1)).
			Return(nil).
			Once()

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
