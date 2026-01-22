package services

import (
	"context"
	"errors"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client_cpf"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	clientFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_Filter(t *testing.T) {
	t.Run("falha quando filtro é nulo", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfFilterService(mockRepo)

		result, err := service.Filter(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha na validação do filtro", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfFilterService(mockRepo)

		invalidFilter := &clientFilter.ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				Limit: -10, // inválido
			},
		}

		result, err := service.Filter(context.Background(), invalidFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfFilterService(mockRepo)

		validFilter := &clientFilter.ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		dbErr := errors.New("falha no banco de dados")

		mockRepo.On("Filter", mock.Anything, validFilter).Return(nil, dbErr).Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso ao retornar lista de clientes", func(t *testing.T) {
		mockRepo := new(mockClient.MockClientCpf)
		service := NewClientCpfFilterService(mockRepo)

		validFilter := &clientFilter.ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockClients := []*model.ClientCpf{
			{ID: 1, Name: "João Silva", Email: *utils.StrToPtr("joao@email.com")},
			{ID: 2, Name: "Maria Souza", Email: *utils.StrToPtr("maria@email.com")},
		}

		mockRepo.On("Filter", mock.Anything, validFilter).Return(mockClients, nil).Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "João Silva", result[0].Name)
		assert.Equal(t, "Maria Souza", result[1].Name)
		mockRepo.AssertExpectations(t)
	})
}
