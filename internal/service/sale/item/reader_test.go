package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	mock_item "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func TestSaleItemService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("id inválido retorna erro", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)

		result, err := service.GetByID(ctx, 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("erro do repositório é propagado", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errors.New("db error"))
		service := NewItemSale(mockRepo)

		result, err := service.GetByID(ctx, 1)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		item := &models.SaleItem{ID: 1}
		mockRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
		service := NewItemSale(mockRepo)

		result, err := service.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, item, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleItemService_GetBySaleID(t *testing.T) {
	ctx := context.Background()

	t.Run("saleID inválido retorna erro", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)

		result, err := service.GetBySaleID(ctx, 0, 10, 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetBySaleID")
	})

	t.Run("paginação inválida retorna erro", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)

		result, err := service.GetBySaleID(ctx, 1, 0, -1)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "GetBySaleID")
	})

	t.Run("erro do repositório é propagado", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		mockRepo.On("GetBySaleID", ctx, int64(1), 10, 0).Return(nil, errors.New("db error"))
		service := NewItemSale(mockRepo)

		result, err := service.GetBySaleID(ctx, 1, 10, 0)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		items := []*models.SaleItem{{ID: 1}, {ID: 2}}
		mockRepo.On("GetBySaleID", ctx, int64(1), 10, 0).Return(items, nil)
		service := NewItemSale(mockRepo)

		result, err := service.GetBySaleID(ctx, 1, 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, items, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleItemService_GetByProductID(t *testing.T) {
	ctx := context.Background()

	t.Run("productID inválido retorna erro", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)

		result, err := service.GetByProductID(ctx, 0, 10, 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetByProductID")
	})

	t.Run("paginação inválida retorna erro", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)

		result, err := service.GetByProductID(ctx, 1, -5, -1)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "GetByProductID")
	})

	t.Run("erro do repositório é propagado", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		mockRepo.On("GetByProductID", ctx, int64(1), 10, 0).Return(nil, errors.New("db error"))
		service := NewItemSale(mockRepo)

		result, err := service.GetByProductID(ctx, 1, 10, 0)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		items := []*models.SaleItem{{ID: 1}, {ID: 2}}
		mockRepo.On("GetByProductID", ctx, int64(1), 10, 0).Return(items, nil)
		service := NewItemSale(mockRepo)

		result, err := service.GetByProductID(ctx, 1, 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, items, result)
		mockRepo.AssertExpectations(t)
	})
}
