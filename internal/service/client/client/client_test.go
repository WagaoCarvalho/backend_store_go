package services

import (
	"context"
	"errors"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/client"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_Create(t *testing.T) {
	t.Run("falha na validação", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		invalidClient := &models.Client{Name: ""} // inválido

		result, err := service.Create(context.Background(), invalidClient)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("erro - cliente duplicado", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

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
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

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
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		cpf := "12345678901"
		client := &models.Client{ID: 1, Name: "Teste", CPF: &cpf}
		expectedErr := errors.New("db error")

		mockRepo.
			On("Create", mock.Anything, client).
			Return(nil, expectedErr)

		result, err := service.Create(context.Background(), client)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, errMsg.ErrCreate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

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
}

func TestClientService_GetByID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("falha ao buscar cliente", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

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
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "teste"}

		mockRepo.
			On("GetByID", mock.Anything, int64(1)).
			Return(client, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, client, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_GetByName(t *testing.T) {
	t.Run("falha por nome vazio", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		result, err := service.GetByName(context.Background(), "")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "GetByName", mock.Anything, mock.Anything)
	})

	t.Run("falha ao buscar cliente por nome", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("GetByName", mock.Anything, "teste").
			Return(nil, expectedErr)

		result, err := service.GetByName(context.Background(), "teste")

		assert.Nil(t, result)
		assert.ErrorContains(t, err, errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "teste"}
		mockRepo.
			On("GetByName", mock.Anything, "teste").
			Return([]*models.Client{client}, nil)

		result, err := service.GetByName(context.Background(), "teste")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, client, result[0])
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_GetVersionByID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		version, err := service.GetVersionByID(context.Background(), 0)

		assert.Zero(t, version)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetVersionByID", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

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
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		mockRepo.
			On("GetVersionByID", mock.Anything, int64(1)).
			Return(5, nil)

		version, err := service.GetVersionByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 5, version)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_GetAll(t *testing.T) {
	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("GetAll", mock.Anything, 10, 0).
			Return(nil, expectedErr)

		result, err := service.GetAll(context.Background(), 10, 0)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste"}
		mockRepo.
			On("GetAll", mock.Anything, 5, 2).
			Return([]*models.Client{client}, nil)

		result, err := service.GetAll(context.Background(), 5, 2)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, client, result[0])
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_ClientExists(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		exists, err := service.ClientExists(context.Background(), 0)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "ClientExists", mock.Anything, mock.Anything)
	})

	t.Run("falha ao verificar existência", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("ClientExists", mock.Anything, int64(1)).
			Return(false, expectedErr)

		exists, err := service.ClientExists(context.Background(), 1)

		assert.False(t, exists)
		assert.ErrorContains(t, err, errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso - cliente existe", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		mockRepo.
			On("ClientExists", mock.Anything, int64(1)).
			Return(true, nil)

		exists, err := service.ClientExists(context.Background(), 1)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_Update(t *testing.T) {
	t.Run("falha - ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 0, Name: "Teste", Version: 1}

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - versão inválida", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 0}

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - validação", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "", Version: 1, CNPJ: nil} // inválido

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha - erro genérico do repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}
		expectedErr := errors.New("db error")

		mockRepo.On("Update", mock.Anything, client).Return(expectedErr).Once()

		err := service.Update(context.Background(), client)
		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		assert.ErrorContains(t, err, expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha - conflito de versão", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}

		mockRepo.On("Update", mock.Anything, client).Return(errMsg.ErrVersionConflict).Once()

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha - cliente duplicado", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}

		mockRepo.On("Update", mock.Anything, client).Return(errMsg.ErrDuplicate).Once()

		err := service.Update(context.Background(), client)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		client := &models.Client{ID: 1, Name: "Teste", Version: 1, CNPJ: utils.StrToPtr("12135135000158")}

		mockRepo.On("Update", mock.Anything, client).Return(nil).Once()

		err := service.Update(context.Background(), client)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_Delete(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("Delete", mock.Anything, int64(1)).
			Return(expectedErr)

		err := service.Delete(context.Background(), 1)

		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		mockRepo.
			On("Delete", mock.Anything, int64(1)).
			Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestClientService_Disable(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		err := service.Disable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Disable", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("Disable", mock.Anything, int64(1)).
			Return(expectedErr)

		err := service.Disable(context.Background(), 1)

		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

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
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		err := service.Enable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Enable", mock.Anything, mock.Anything)
	})

	t.Run("falha no repo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		expectedErr := errors.New("db error")
		mockRepo.
			On("Enable", mock.Anything, int64(1)).
			Return(expectedErr)

		err := service.Enable(context.Background(), 1)

		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientRepository)
		service := NewClient(mockRepo)

		mockRepo.
			On("Enable", mock.Anything, int64(1)).
			Return(nil)

		err := service.Enable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
