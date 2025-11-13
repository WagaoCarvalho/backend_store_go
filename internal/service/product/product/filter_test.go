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

func TestProductService_Filter(t *testing.T) {

	setup := func() (*mockProduct.ProductMock, ProductService) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: filterData é nil", func(t *testing.T) {
		mockRepo, service := setup()

		products, err := service.Filter(context.Background(), nil) // Captura ambos os retornos

		assert.Nil(t, products)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Filter")
	})

	t.Run("falha: validação do filtro", func(t *testing.T) {
		mockRepo, service := setup()

		invalidFilter := &models.ProductFilter{
			MinCostPrice: utils.Float64Ptr(100.0),
			MaxCostPrice: utils.Float64Ptr(50.0), // Min > Max - inválido
		}

		products, err := service.Filter(context.Background(), invalidFilter) // Captura ambos os retornos

		assert.Nil(t, products)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de preço de custo inválido")
		mockRepo.AssertNotCalled(t, "Filter")
	})

	t.Run("sucesso: filtro aplicado", func(t *testing.T) {
		mockRepo, service := setup()
		ctx := context.Background()

		filterData := &models.ProductFilter{
			ProductName: "Notebook",
			Status:      utils.BoolPtr(true),
		}

		expectedProducts := []*models.Product{
			{ID: 1, ProductName: "Notebook Dell", Status: true},
			{ID: 2, ProductName: "Notebook HP", Status: true},
		}

		mockRepo.
			On("Filter", ctx, filterData).
			Return(expectedProducts, nil)

		result, err := service.Filter(ctx, filterData)

		assert.NoError(t, err)
		assert.Len(t, result, len(expectedProducts))
		assert.Equal(t, expectedProducts, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso: filtro vazio retorna todos os produtos", func(t *testing.T) {
		mockRepo, service := setup()
		ctx := context.Background()

		filterData := &models.ProductFilter{} // Filtro vazio

		expectedProducts := []*models.Product{
			{ID: 1, ProductName: "Produto 1"},
			{ID: 2, ProductName: "Produto 2"},
			{ID: 3, ProductName: "Produto 3"},
		}

		mockRepo.
			On("Filter", ctx, filterData).
			Return(expectedProducts, nil)

		result, err := service.Filter(ctx, filterData)

		assert.NoError(t, err)
		assert.Len(t, result, len(expectedProducts))
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo, service := setup()
		ctx := context.Background()

		filterData := &models.ProductFilter{
			ProductName: "Test",
		}

		mockErr := errors.New("erro no banco de dados")
		mockRepo.
			On("Filter", ctx, filterData).
			Return(nil, mockErr)

		result, err := service.Filter(ctx, filterData)

		assert.Nil(t, result)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso: nenhum produto encontrado", func(t *testing.T) {
		mockRepo, service := setup()
		ctx := context.Background()

		filterData := &models.ProductFilter{
			ProductName: "ProdutoInexistente",
		}

		emptyProducts := []*models.Product{}

		mockRepo.
			On("Filter", ctx, filterData).
			Return(emptyProducts, nil)

		result, err := service.Filter(ctx, filterData)

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}
