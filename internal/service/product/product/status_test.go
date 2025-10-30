package services

import (
	"context"
	"fmt"
	"testing"

	mock_product "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_DisableProduct(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProduct(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		invalidID := int64(0)
		err := service.DisableProduct(context.Background(), invalidID)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "DisableProduct")
	})

	t.Run("Deve desabilitar produto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DisableProduct", mock.Anything, int64(1)).Return(nil).Once()

		err := service.DisableProduct(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao desabilitar produto", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DisableProduct", mock.Anything, int64(2)).Return(fmt.Errorf("erro banco")).Once()

		err := service.DisableProduct(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao desabilitar")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_EnableProduct(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProduct(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		invalidID := int64(0)
		err := service.EnableProduct(context.Background(), invalidID)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "EnableProduct")
	})

	t.Run("Deve habilitar produto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableProduct", mock.Anything, int64(1)).Return(nil).Once()

		err := service.EnableProduct(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao habilitar produto", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableProduct", mock.Anything, int64(2)).Return(fmt.Errorf("erro banco")).Once()

		err := service.EnableProduct(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao ativar")
		mockRepo.AssertExpectations(t)
	})
}
