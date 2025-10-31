package services

import (
	"context"
	"errors"
	"testing"

	mockSale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
)

func TestSaleService_Cancel(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSale(mockRepo)
	ctx := context.Background()

	t.Run("id inválido", func(t *testing.T) {
		err := svc.Cancel(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("sale não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errMsg.ErrNotFound).Once()
		err := svc.Cancel(ctx, 1)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale não ativa", func(t *testing.T) {
		sale := &models.Sale{ID: 2, Status: "completed"}
		mockRepo.On("GetByID", ctx, int64(2)).Return(sale, nil).Once()
		err := svc.Cancel(ctx, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only active sales can be canceled")
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo update", func(t *testing.T) {
		sale := &models.Sale{ID: 3, Status: "active"}
		mockRepo.On("GetByID", ctx, int64(3)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(errors.New("update error")).Once()

		err := svc.Cancel(ctx, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		sale := &models.Sale{ID: 4, Status: "active"}
		mockRepo.On("GetByID", ctx, int64(4)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(nil).Once()

		err := svc.Cancel(ctx, 4)
		assert.NoError(t, err)
		assert.Equal(t, "canceled", sale.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo get", func(t *testing.T) {
		mockRepo := new(mockSale.MockSale)
		svc := NewSale(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, int64(5)).Return(nil, errors.New("db error")).Once()

		err := svc.Cancel(ctx, 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

}

func TestSaleService_Complete(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSale(mockRepo)
	ctx := context.Background()

	t.Run("id inválido", func(t *testing.T) {
		err := svc.Complete(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("sale não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errMsg.ErrNotFound).Once()
		err := svc.Complete(ctx, 1)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale não ativa", func(t *testing.T) {
		sale := &models.Sale{ID: 2, Status: "canceled"}
		mockRepo.On("GetByID", ctx, int64(2)).Return(sale, nil).Once()
		err := svc.Complete(ctx, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only active sales can be completed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo update", func(t *testing.T) {
		sale := &models.Sale{ID: 3, Status: "active"}
		mockRepo.On("GetByID", ctx, int64(3)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(errors.New("update error")).Once()

		err := svc.Complete(ctx, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		sale := &models.Sale{ID: 4, Status: "active"}
		mockRepo.On("GetByID", ctx, int64(4)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(nil).Once()

		err := svc.Complete(ctx, 4)
		assert.NoError(t, err)
		assert.Equal(t, "completed", sale.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo get genérico", func(t *testing.T) {
		mockRepo := new(mockSale.MockSale)
		svc := NewSale(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, int64(10)).Return(nil, errors.New("erro de banco")).Once()

		err := svc.Complete(ctx, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

}
