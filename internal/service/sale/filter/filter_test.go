package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockSale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	saleFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleService_Filter(t *testing.T) {

	t.Run("falha quando filtro é nulo", func(t *testing.T) {
		mockRepo := new(mockSale.MockSale)
		service := NewSaleFilterService(mockRepo)

		result, err := service.Filter(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha na validação do filtro", func(t *testing.T) {
		mockRepo := new(mockSale.MockSale)
		service := NewSaleFilterService(mockRepo)

		invalidFilter := &saleFilter.SaleFilter{
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
		mockRepo := new(mockSale.MockSale)
		service := NewSaleFilterService(mockRepo)

		validFilter := &saleFilter.SaleFilter{
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

	t.Run("sucesso ao retornar lista de vendas", func(t *testing.T) {
		mockRepo := new(mockSale.MockSale)
		service := NewSaleFilterService(mockRepo)

		validFilter := &saleFilter.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		now := time.Now()

		mockSales := []*model.Sale{
			{
				ID:          1,
				SaleDate:    now,
				TotalAmount: 100,
				PaymentType: "cash",
				Status:      "completed",
				Version:     1,
			},
			{
				ID:          2,
				SaleDate:    now,
				TotalAmount: 250,
				PaymentType: "card",
				Status:      "active",
				Version:     1,
			},
		}

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(mockSales, nil).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "cash", result[0].PaymentType)
		assert.Equal(t, int64(2), result[1].ID)
		assert.Equal(t, "card", result[1].PaymentType)

		mockRepo.AssertExpectations(t)
	})
}
