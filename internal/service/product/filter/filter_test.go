package services

import (
	"context"
	"errors"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	filterProduct "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_Filter(t *testing.T) {

	t.Run("falha quando filtro é nulo", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductFilterService(mockRepo)

		result, err := service.Filter(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha na validação do filtro", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductFilterService(mockRepo)

		invalidFilter := &filterProduct.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit: -1, // inválido
			},
		}

		result, err := service.Filter(context.Background(), invalidFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductFilterService(mockRepo)

		validFilter := &filterProduct.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		dbErr := errors.New("falha no banco de dados")

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(nil, dbErr).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso ao retornar lista de produtos", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductFilterService(mockRepo)

		validFilter := &filterProduct.ProductFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockProducts := []*model.Product{
			{
				ID:           1,
				ProductName:  "Produto A",
				Manufacturer: "Fabricante X",
				Status:       true,
				Version:      1,
			},
			{
				ID:           2,
				ProductName:  "Produto B",
				Manufacturer: "Fabricante Y",
				Status:       true,
				Version:      1,
			},
		}

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(mockProducts, nil).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Produto A", result[0].ProductName)
		assert.Equal(t, "Produto B", result[1].ProductName)
		mockRepo.AssertExpectations(t)
	})
}
