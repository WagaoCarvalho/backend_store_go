package services

import (
	"context"
	"errors"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_Create(t *testing.T) {
	t.Run("falha na validação", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		invalidClient := &models.Client{Name: ""} // inválido

		result, err := service.Create(context.Background(), invalidClient)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("erro - cliente duplicado", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{
			Name:  "Cliente Teste",
			Email: utils.StrToPtr("teste@teste.com"),
			CNPJ:  utils.StrToPtr("12345678000199"),
		}

		mockRepo.On("Create", mock.Anything, client).Return(nil, errMsg.ErrDuplicate).Once()

		result, err := service.Create(context.Background(), client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - chave estrangeira inválida", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{
			Name:  "Cliente Teste",
			Email: utils.StrToPtr("teste@teste.com"),
			CNPJ:  utils.StrToPtr("12345678000199"),
		}

		mockRepo.On("Create", mock.Anything, client).Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		result, err := service.Create(context.Background(), client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao criar no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		cpf := "12345678901"
		client := &models.Client{ID: 1, Name: "Teste", CPF: &cpf}
		expectedErr := errors.New("erro ao criar")

		mockRepo.
			On("Create", mock.Anything, client).
			Return(nil, expectedErr)

		result, err := service.Create(context.Background(), client)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, errMsg.ErrCreate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		cpf := "12345678901"
		client := &models.Client{ID: 1, Name: "Teste", CPF: &cpf}

		mockRepo.
			On("Create", mock.Anything, client).
			Return(client, nil)

		result, err := service.Create(context.Background(), client)

		assert.NoError(t, err)
		assert.Equal(t, client, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha quando client é nil", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		result, err := service.Create(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create")
	})

}

func TestClientService_Update(t *testing.T) {
	t.Run("falha - ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{ID: 0, Name: "Teste", Version: 1}

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - versão inválida", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 0}

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrZeroVersion)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - validação", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{ID: 1, Name: "", Version: 1, CNPJ: nil} // inválido

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - erro genérico do repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}
		expectedErr := errors.New("db error")

		mockRepo.On("Update", mock.Anything, client).Return(expectedErr).Once()

		err := service.Update(context.Background(), client)
		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		assert.ErrorContains(t, err, expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha - conflito de versão", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}

		mockRepo.On("Update", mock.Anything, client).Return(errMsg.ErrZeroVersion).Once()

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrZeroVersion)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha - cliente duplicado", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}

		mockRepo.On("Update", mock.Anything, client).Return(errMsg.ErrDuplicate).Once()

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}

		mockRepo.On("Update", mock.Anything, client).Return(nil).Once()

		err := service.Update(context.Background(), client)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_Delete(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("Delete", mock.Anything, int64(1)).
			Return(expectedErr)

		err := service.Delete(context.Background(), 1)

		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClient)
		service := NewClientService(mockRepo)

		mockRepo.
			On("Delete", mock.Anything, int64(1)).
			Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
