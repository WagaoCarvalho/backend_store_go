package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_EnableDiscount(t *testing.T) {

	setup := func() (*mockProduct.ProductMock, ProductService) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)
		return mockRepo, service
	}

	t.Run("Deve habilitar desconto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableDiscount", mock.Anything, int64(1)).Return(nil).Once()

		err := service.EnableDiscount(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)
		service := product{repo: repoMock}

		ctx := context.Background() // ⚠️ necessário criar o contexto
		err := service.EnableDiscount(ctx, 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "EnableDiscount")
	})

	t.Run("Deve retornar erro quando repo falhar", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableDiscount", mock.Anything, int64(1)).Return(errMsg.ErrNotFound).Once()

		err := service.EnableDiscount(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrProductEnableDiscount)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_DisableDiscount(t *testing.T) {

	setup := func() (*mockProduct.ProductMock, ProductService) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)
		return mockRepo, service
	}

	t.Run("Deve desabilitar desconto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DisableDiscount", mock.Anything, int64(1)).Return(nil).Once()

		err := service.DisableDiscount(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)
		service := product{repo: repoMock}

		ctx := context.Background() // ⚠️ necessário criar o contexto
		err := service.DisableDiscount(ctx, 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "DisableDiscount")
	})

	t.Run("Erro: produto não encontrado", func(t *testing.T) {
		mockRepo, service := setup()
		mockRepo.On("EnableDiscount", mock.Anything, int64(2)).Return(errMsg.ErrNotFound).Once()

		err := service.EnableDiscount(context.Background(), 2)

		assert.ErrorIs(t, err, errMsg.ErrProductEnableDiscount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro inesperado do repositório", func(t *testing.T) {
		mockRepo, service := setup()
		unexpectedErr := fmt.Errorf("erro de conexão")
		mockRepo.On("DisableDiscount", mock.Anything, int64(3)).Return(unexpectedErr).Once()

		err := service.DisableDiscount(context.Background(), 3)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrProductDisableDiscount)
		assert.Contains(t, err.Error(), "erro de conexão")
		mockRepo.AssertExpectations(t)
	})

}

func TestProductService_ApplyDiscount(t *testing.T) {

	setup := func() (*mockProduct.ProductMock, ProductService) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		ctx := context.Background()
		product, err := service.ApplyDiscount(ctx, 0, 10.0) // ID inválido

		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "ApplyDiscount")
	})

	t.Run("falha: percent inválido", func(t *testing.T) {
		mockRepo, service := setup()

		ctx := context.Background()
		product, err := service.ApplyDiscount(ctx, 1, 0) // percent inválido

		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrPercentInvalid)
		mockRepo.AssertNotCalled(t, "ApplyDiscount")
	})

	t.Run("Deve aplicar desconto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(1)
		percent := 10.0

		expectedProduct := &models.Product{
			ID:            productID,
			ProductName:   "Produto Teste",
			SalePrice:     90.0,
			AllowDiscount: true,
		}

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(expectedProduct, nil).
			Once()

		product, err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, expectedProduct, product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro se produto não encontrado", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(99)
		percent := 15.0

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil, errMsg.ErrNotFound).
			Once()

		product, err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrProductApplyDiscount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico ao aplicar desconto", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(2)
		percent := 5.0
		expectedErr := errors.New("erro inesperado no banco")

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil, expectedErr).
			Once()

		product, err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrProductApplyDiscount)
		mockRepo.AssertExpectations(t)
	})
}
