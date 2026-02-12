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

func TestProductService_EnableDiscount(t *testing.T) {
	setup := func() (*mockProduct.ProductMock, ProductService) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductService(mockRepo)
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
		service := productService{repo: repoMock}

		ctx := context.Background()
		err := service.EnableDiscount(ctx, 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "EnableDiscount")
	})

	t.Run("Deve propagar erro NotFound do repositório", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableDiscount", mock.Anything, int64(1)).Return(errMsg.ErrNotFound).Once()

		err := service.EnableDiscount(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrNotFound) // Agora propaga, não envolve
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve envolver erro genérico do repositório", func(t *testing.T) {
		mockRepo, service := setup()

		repoErr := fmt.Errorf("erro de conexão")
		mockRepo.On("EnableDiscount", mock.Anything, int64(1)).Return(repoErr).Once()

		err := service.EnableDiscount(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrProductEnableDiscount)
		assert.Contains(t, err.Error(), "erro de conexão")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_DisableDiscount(t *testing.T) {
	setup := func() (*mockProduct.ProductMock, ProductService) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductService(mockRepo)
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
		service := productService{repo: repoMock}

		ctx := context.Background()
		err := service.DisableDiscount(ctx, 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		repoMock.AssertNotCalled(t, "DisableDiscount")
	})

	t.Run("Deve propagar erro NotFound do repositório", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DisableDiscount", mock.Anything, int64(2)).Return(errMsg.ErrNotFound).Once()

		err := service.DisableDiscount(context.Background(), 2)

		assert.ErrorIs(t, err, errMsg.ErrNotFound) // Propaga, não envolve
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve envolver erro genérico do repositório", func(t *testing.T) {
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
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		ctx := context.Background()
		err := service.ApplyDiscount(ctx, 0, 10.0) // ID inválido

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "ApplyDiscount")
	})

	t.Run("falha: percent negativo", func(t *testing.T) {
		mockRepo, service := setup()

		ctx := context.Background()
		err := service.ApplyDiscount(ctx, 1, -5.0)

		assert.ErrorIs(t, err, errMsg.ErrInvalidDiscountPercent)
		mockRepo.AssertNotCalled(t, "ApplyDiscount")
	})

	t.Run("falha: percent maior que 100", func(t *testing.T) {
		mockRepo, service := setup()

		ctx := context.Background()
		err := service.ApplyDiscount(ctx, 1, 150.0)

		assert.ErrorIs(t, err, errMsg.ErrInvalidDiscountPercent)
		mockRepo.AssertNotCalled(t, "ApplyDiscount")
	})

	t.Run("Deve aplicar desconto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(1)
		percent := 10.0

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil).
			Once()

		err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve propagar erro NotFound do repositório", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(99)
		percent := 15.0

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(errMsg.ErrNotFound).
			Once()

		err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.ErrorIs(t, err, errMsg.ErrNotFound) // Propaga, não envolve
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve propagar erro DiscountNotAllowed do repositório", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(2)
		percent := 5.0

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(errMsg.ErrProductDiscountNotAllowed).
			Once()

		err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.ErrorIs(t, err, errMsg.ErrProductDiscountNotAllowed) // Propaga
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve envolver erro genérico do repositório", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(3)
		percent := 20.0
		expectedErr := fmt.Errorf("erro inesperado no banco")

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(expectedErr).
			Once()

		err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrProductApplyDiscount)
		assert.Contains(t, err.Error(), "erro inesperado no banco")
		mockRepo.AssertExpectations(t)
	})
}
