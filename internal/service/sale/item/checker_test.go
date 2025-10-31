package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	mock_item "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func TestSaleItemChecker_ItemExists(t *testing.T) {
	ctx := context.Background()

	t.Run("id inválido retorna erro", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)

		exists, err := service.ItemExists(ctx, 0)
		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "ItemExists")
	})

	t.Run("item existe", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)
		mockRepo.On("ItemExists", ctx, int64(10)).Return(true, nil)

		exists, err := service.ItemExists(ctx, 10)
		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("item não existe", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)
		mockRepo.On("ItemExists", ctx, int64(99)).Return(false, nil)

		exists, err := service.ItemExists(ctx, 99)
		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mock_item.MockSaleItem)
		service := NewItemSale(mockRepo)
		mockRepo.On("ItemExists", ctx, int64(5)).Return(false, errors.New("db error"))

		exists, err := service.ItemExists(ctx, 5)
		assert.False(t, exists)
		assert.ErrorContains(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}
