package services

import (
	"context"
	"errors"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestProductService_Create(t *testing.T) {
	ctx := context.Background()

	validProduct := func() *models.Product {
		return &models.Product{
			ProductName:   "Produto Teste",
			Manufacturer:  "Fabricante X",
			SupplierID:    utils.Int64Ptr(1),
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 5,
			Status:        true,
		}
	}

	t.Run("falha: produto é nil", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		// Chama com produto nulo
		created, err := service.Create(ctx, nil)

		assert.Nil(t, created)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

		service := NewProduct(mockRepo)

		input := validProduct()

		mockRepo.On("Create", ctx, input).Return(&models.Product{
			ID:            1,
			ProductName:   input.ProductName,
			Manufacturer:  input.Manufacturer,
			SupplierID:    input.SupplierID,
			CostPrice:     input.CostPrice,
			SalePrice:     input.SalePrice,
			StockQuantity: input.StockQuantity,
			Status:        input.Status,
		}, nil).Once()

		created, err := service.Create(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.Equal(t, int64(1), created.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("produto inválido", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := &models.Product{
			Status: true, // faltando campos obrigatórios
		}

		created, err := service.Create(ctx, input)

		assert.Nil(t, created)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

		service := NewProduct(mockRepo)

		input := validProduct()

		mockErr := errors.New("erro no repositório")
		mockRepo.On("Create", ctx, input).Return(nil, mockErr).Once()

		created, err := service.Create(ctx, input)

		assert.Nil(t, created)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Update(t *testing.T) {
	ctx := context.Background()

	validProduct := func() *models.Product {
		return &models.Product{
			ID:            1,
			ProductName:   "Produto Atualizado",
			Manufacturer:  "Fabricante X",
			SupplierID:    utils.Int64Ptr(1),
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 5,
			Status:        true,
			Version:       1,
		}
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()
		input.ID = 0

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: validação inválida", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()
		input.ProductName = "" // invalida a validação

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: versão inválida", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()
		input.Version = 0

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: produto não encontrado (ProductExists = false)", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()

		// Simula o update retornando ErrNotFound
		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		// Simula o ProductExists retornando false
		mockRepo.On("ProductExists", ctx, input.ID).Return(false, nil).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: conflito de versão (ProductExists = true)", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("ProductExists", ctx, input.ID).Return(true, nil).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro ao verificar existência", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()
		expectedErr := errors.New("erro no banco")

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("ProductExists", ctx, input.ID).Return(false, expectedErr).Once()

		err := service.Update(ctx, input)

		assert.ErrorContains(t, err, errMsg.ErrGet.Error())
		assert.ErrorContains(t, err, expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro genérico no Update", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()
		expectedErr := errors.New("erro genérico")

		mockRepo.On("Update", ctx, input).Return(expectedErr).Once()

		err := service.Update(ctx, input)

		assert.ErrorContains(t, err, errMsg.ErrUpdate.Error())
		assert.ErrorContains(t, err, expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProduct(mockRepo)

		input := validProduct()

		mockRepo.On("Update", ctx, input).Return(nil).Once()

		err := service.Update(ctx, input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

		service := NewProduct(mockRepo)

		id := int64(1)

		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

		service := NewProduct(mockRepo)

		invalidID := int64(0)
		err := service.Delete(context.Background(), invalidID)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

		service := NewProduct(mockRepo)

		id := int64(1)
		mockErr := errors.New("erro no repositório")

		mockRepo.On("Delete", ctx, id).Return(mockErr)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}
