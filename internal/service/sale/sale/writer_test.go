package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockSale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
)

func TestSaleService_Create(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSale(mockRepo)
	ctx := context.Background()

	t.Run("sale nil", func(t *testing.T) {
		result, err := svc.Create(ctx, nil)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("structural validation fails", func(t *testing.T) {
		s := &models.Sale{
			TotalAmount:   -1,
			TotalDiscount: -1,
			PaymentType:   "",
			Status:        "",
		}
		result, err := svc.Create(ctx, s)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("business validation fails", func(t *testing.T) {
		s := &models.Sale{
			TotalAmount:   10,
			TotalDiscount: 20, // maior que o total
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
		}
		result, err := svc.Create(ctx, s)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repo returns error", func(t *testing.T) {
		s := &models.Sale{
			TotalAmount:   100,
			TotalDiscount: 10,
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
		}
		mockRepo.On("Create", ctx, s).Return(nil, errors.New("repo error")).Once()

		result, err := svc.Create(ctx, s)
		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		s := &models.Sale{
			TotalAmount:   100,
			TotalDiscount: 10,
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
			Notes:         "Teste",
		}
		created := &models.Sale{
			ID:            1,
			TotalAmount:   s.TotalAmount,
			TotalDiscount: s.TotalDiscount,
			PaymentType:   s.PaymentType,
			Status:        s.Status,
			SaleDate:      s.SaleDate,
			Notes:         s.Notes,
		}
		mockRepo.On("Create", ctx, s).Return(created, nil).Once()

		result, err := svc.Create(ctx, s)
		assert.NoError(t, err)
		assert.Equal(t, created, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_Update(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSale(mockRepo)
	ctx := context.Background()

	t.Run("sale nil", func(t *testing.T) {
		err := svc.Update(ctx, nil)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("id zero", func(t *testing.T) {
		s := &models.Sale{
			ID:            0,
			Version:       1,
			TotalAmount:   100,
			TotalDiscount: 10,
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("version zero", func(t *testing.T) {
		s := &models.Sale{
			ID:            1,
			Version:       0,
			TotalAmount:   100,
			TotalDiscount: 10,
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
	})

	t.Run("structural validation fails", func(t *testing.T) {
		s := &models.Sale{
			ID:            1,
			Version:       1,
			TotalAmount:   -1,
			TotalDiscount: -1,
			PaymentType:   "",
			Status:        "",
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("business validation fails", func(t *testing.T) {
		s := &models.Sale{
			ID:            1,
			Version:       1,
			TotalAmount:   10,
			TotalDiscount: 20, // maior que o total
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repo returns error", func(t *testing.T) {
		s := &models.Sale{
			ID:            1,
			Version:       1,
			TotalAmount:   100,
			TotalDiscount: 10,
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
		}
		mockRepo.On("Update", ctx, s).Return(errors.New("repo error")).Once()

		err := svc.Update(ctx, s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		s := &models.Sale{
			ID:            1,
			Version:       1,
			TotalAmount:   100,
			TotalDiscount: 10,
			PaymentType:   "cash",
			Status:        "active",
			SaleDate:      time.Now(),
			Notes:         "Teste",
		}
		mockRepo.On("Update", ctx, s).Return(nil).Once()

		err := svc.Update(ctx, s)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_Delete(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSale(mockRepo)
	ctx := context.Background()

	t.Run("id zero", func(t *testing.T) {
		err := svc.Delete(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("repo returns error", func(t *testing.T) {
		mockRepo.On("Delete", ctx, int64(1)).Return(errors.New("repo error")).Once()
		err := svc.Delete(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrDelete.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, int64(2)).Return(nil).Once()
		err := svc.Delete(ctx, 2)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
