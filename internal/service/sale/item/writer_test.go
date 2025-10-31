package services

import (
	"context"
	"errors"
	"testing"

	mockItem "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
)

func TestSaleItemService_Create(t *testing.T) {
	mockRepo := new(mockItem.MockSaleItem)
	svc := NewItemSale(mockRepo)
	ctx := context.Background()

	t.Run("item nil", func(t *testing.T) {
		result, err := svc.Create(ctx, nil)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("structural validation fails", func(t *testing.T) {
		i := &models.SaleItem{
			SaleID:    0,
			ProductID: -1,
			Quantity:  0,
			UnitPrice: -10,
			Discount:  -1,
			Tax:       -1,
			Subtotal:  -1,
		}
		result, err := svc.Create(ctx, i)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("business validation fails", func(t *testing.T) {
		i := &models.SaleItem{
			SaleID:    1,
			ProductID: 1,
			Quantity:  2,
			UnitPrice: 10,
			Discount:  1,
			Tax:       0,
			Subtotal:  10, // errado: deveria ser 19
		}
		result, err := svc.Create(ctx, i)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repo returns error", func(t *testing.T) {
		i := &models.SaleItem{
			SaleID:    1,
			ProductID: 1,
			Quantity:  2,
			UnitPrice: 10,
			Discount:  1,
			Tax:       0,
			Subtotal:  19, // correto: 2*10 - 1 + 0
		}

		mockRepo.On("Create", ctx, i).Return(nil, errors.New("repo error")).Once()

		result, err := svc.Create(ctx, i)
		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		i := &models.SaleItem{
			SaleID:    1,
			ProductID: 1,
			Quantity:  2,
			UnitPrice: 10,
			Discount:  1,
			Tax:       0,
			Subtotal:  19, // válido
		}

		created := &models.SaleItem{
			ID:        1,
			SaleID:    i.SaleID,
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
			UnitPrice: i.UnitPrice,
			Discount:  i.Discount,
			Tax:       i.Tax,
			Subtotal:  i.Subtotal,
		}

		mockRepo.On("Create", ctx, i).Return(created, nil).Once()

		result, err := svc.Create(ctx, i)
		assert.NoError(t, err)
		assert.Equal(t, created, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleItemService_Update(t *testing.T) {
	mockRepo := new(mockItem.MockSaleItem)
	svc := NewItemSale(mockRepo)
	ctx := context.Background()

	t.Run("item nil", func(t *testing.T) {
		err := svc.Update(ctx, nil)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("structural validation fails", func(t *testing.T) {
		i := &models.SaleItem{
			SaleID:    0,
			ProductID: -1,
			Quantity:  0,
			UnitPrice: -10,
			Discount:  -1,
			Tax:       -1,
			Subtotal:  -1,
		}
		err := svc.Update(ctx, i)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("business validation fails", func(t *testing.T) {
		i := &models.SaleItem{
			SaleID:    1,
			ProductID: 1,
			Quantity:  2,
			UnitPrice: 10,
			Discount:  1,
			Tax:       0,
			Subtotal:  0, // incorreto
		}
		err := svc.Update(ctx, i)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repo returns error", func(t *testing.T) {
		i := &models.SaleItem{
			ID:        1,
			SaleID:    1,
			ProductID: 1,
			Quantity:  2,
			UnitPrice: 10,
			Discount:  1,
			Tax:       0,
			Subtotal:  19,
		}

		mockRepo.On("Update", ctx, i).Return(errors.New("repo error")).Once()

		err := svc.Update(ctx, i)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		i := &models.SaleItem{
			ID:        1,
			SaleID:    1,
			ProductID: 1,
			Quantity:  2,
			UnitPrice: 10,
			Discount:  1,
			Tax:       0,
			Subtotal:  19,
		}

		mockRepo.On("Update", ctx, i).Return(nil).Once()

		err := svc.Update(ctx, i)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleItemService_Delete(t *testing.T) {
	mockRepo := new(mockItem.MockSaleItem)
	svc := NewItemSale(mockRepo)
	ctx := context.Background()

	t.Run("id zero", func(t *testing.T) {
		err := svc.Delete(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("repo returns error", func(t *testing.T) {
		mockRepo.On("Delete", ctx, int64(1)).Return(errors.New("repo error")).Once()
		err := svc.Delete(ctx, 1)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, int64(2)).Return(nil).Once()
		err := svc.Delete(ctx, 2)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleItemService_DeleteBySaleID(t *testing.T) {
	mockRepo := new(mockItem.MockSaleItem)
	svc := NewItemSale(mockRepo)
	ctx := context.Background()

	t.Run("saleID zero", func(t *testing.T) {
		err := svc.DeleteBySaleID(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("repo returns error", func(t *testing.T) {
		mockRepo.On("DeleteBySaleID", ctx, int64(1)).Return(errors.New("repo error")).Once()
		err := svc.DeleteBySaleID(ctx, 1)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteBySaleID", ctx, int64(2)).Return(nil).Once()
		err := svc.DeleteBySaleID(ctx, 2)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
