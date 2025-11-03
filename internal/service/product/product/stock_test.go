package services

import (
	"context"
	"fmt"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_UpdateStock(t *testing.T) {

	setup := func() (*mockProduct.ProductMock, ProductService) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)
		return mockRepo, service
	}

	t.Run("Deve atualizar o estoque com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("UpdateStock", mock.Anything, int64(1), 25).Return(nil).Once()

		err := service.UpdateStock(context.Background(), 1, 25)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		err := service.UpdateStock(context.Background(), 0, 10)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "UpdateStock")
	})

	t.Run("falha: quantidade inválida", func(t *testing.T) {
		mockRepo, service := setup()

		err := service.UpdateStock(context.Background(), 1, 0)

		assert.ErrorIs(t, err, errMsg.ErrInvalidQuantity)
		mockRepo.AssertNotCalled(t, "UpdateStock")
	})

	t.Run("Deve retornar erro quando repo falhar", func(t *testing.T) {
		mockRepo, service := setup()
		expectedErr := fmt.Errorf("erro de banco")

		mockRepo.On("UpdateStock", mock.Anything, int64(1), 25).Return(expectedErr).Once()

		err := service.UpdateStock(context.Background(), 1, 25)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_IncreaseStock(t *testing.T) {
	ctx := context.Background()

	t.Run("Deve retornar erro quando repo retornar erro", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)

		service := product{repo: repoMock}

		repoMock.On("IncreaseStock", ctx, int64(1), 10).Return(errMsg.ErrNotFound)

		err := service.IncreaseStock(ctx, 1, 10)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		repoMock.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)
		service := product{repo: repoMock}

		err := service.IncreaseStock(ctx, 0, 10)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "IncreaseStock")
	})

	t.Run("falha: quantidade inválida", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)
		service := product{repo: repoMock}

		err := service.IncreaseStock(ctx, 1, 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "IncreaseStock")
	})

	t.Run("Deve aumentar estoque com sucesso", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)

		service := product{repo: repoMock}

		repoMock.On("IncreaseStock", ctx, int64(1), 5).Return(nil)

		err := service.IncreaseStock(ctx, 1, 5)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestProductService_DecreaseStock(t *testing.T) {
	ctx := context.Background()

	t.Run("Deve retornar erro quando repo retornar erro", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)

		service := product{repo: repoMock}

		repoMock.On("DecreaseStock", ctx, int64(1), 10).Return(errMsg.ErrNotFound)

		err := service.DecreaseStock(ctx, 1, 10)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		repoMock.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)
		service := product{repo: repoMock}

		err := service.DecreaseStock(ctx, 0, 10)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "DecreaseStock")
	})

	t.Run("falha: quantidade inválida", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)
		service := product{repo: repoMock}

		err := service.DecreaseStock(ctx, 1, 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "DecreaseStock")
	})

	t.Run("Deve diminuir estoque com sucesso", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)

		service := product{repo: repoMock}

		repoMock.On("DecreaseStock", ctx, int64(1), 10).Return(nil)

		err := service.DecreaseStock(ctx, 1, 10)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestProductService_GetStock(t *testing.T) {
	ctx := context.Background()

	t.Run("Deve retornar erro quando repo retornar erro", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)

		service := product{repo: repoMock}

		repoMock.On("GetStock", ctx, int64(1)).Return(0, fmt.Errorf("erro inesperado"))

		stock, err := service.GetStock(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, 0, stock)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		repoMock.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)
		service := product{repo: repoMock}

		stock, err := service.GetStock(ctx, 0)

		assert.Equal(t, 0, stock)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "GetStock")
	})

	t.Run("Deve retornar estoque com sucesso", func(t *testing.T) {
		repoMock := new(mockProduct.ProductMock)

		service := product{repo: repoMock}

		repoMock.On("GetStock", ctx, int64(1)).Return(25, nil)

		stock, err := service.GetStock(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 25, stock)
		repoMock.AssertExpectations(t)
	})
}
