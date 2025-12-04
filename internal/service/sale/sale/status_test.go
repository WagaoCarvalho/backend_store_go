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

func TestSaleService_GetByStatus(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()

	t.Run("status inválido", func(t *testing.T) {
		result, err := svc.GetByStatus(ctx, "", 10, 0, "id", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("erro paginação", func(t *testing.T) {
		result, err := svc.GetByStatus(ctx, "active", 0, 0, "id", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidLimit)
	})

	t.Run("erro order field", func(t *testing.T) {
		result, err := svc.GetByStatus(ctx, "active", 10, 0, "invalid", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderField)
	})

	t.Run("erro order direction", func(t *testing.T) {
		result, err := svc.GetByStatus(ctx, "active", 10, 0, "id", "invalid")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderDirection)
	})

	t.Run("erro genérico do repo", func(t *testing.T) {
		expectedErr := errors.New("repo error")
		mockRepo.On("GetByStatus", ctx, "active", 10, 0, "id", "asc").Return(nil, expectedErr).Once()
		result, err := svc.GetByStatus(ctx, "active", 10, 0, "id", "asc")
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		expectedSales := []*models.Sale{{ID: 1}, {ID: 2}}
		mockRepo.On("GetByStatus", ctx, "active", 10, 0, "id", "asc").Return(expectedSales, nil).Once()
		result, err := svc.GetByStatus(ctx, "active", 10, 0, "id", "asc")
		assert.NoError(t, err)
		assert.Equal(t, expectedSales, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_Cancel(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
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
		assert.Contains(t, err.Error(), "somente vendas ativas podem ser canceladas")
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
		svc := NewSaleService(mockRepo)
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
	svc := NewSaleService(mockRepo)
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
		assert.Contains(t, err.Error(), "somente vendas ativas podem ser concluídas")
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
		svc := NewSaleService(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, int64(10)).Return(nil, errors.New("erro de banco")).Once()

		err := svc.Complete(ctx, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

}

func TestSaleService_Returned(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()

	t.Run("id inválido", func(t *testing.T) {
		err := svc.Returned(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("sale não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errMsg.ErrNotFound).Once()
		err := svc.Returned(ctx, 1)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo get", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(2)).Return(nil, errors.New("db error")).Once()

		err := svc.Returned(ctx, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale não concluída", func(t *testing.T) {
		sale := &models.Sale{ID: 3, Status: "active"}
		mockRepo.On("GetByID", ctx, int64(3)).Return(sale, nil).Once()
		err := svc.Returned(ctx, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "somente vendas concluídas podem ser devolvidas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale status pending", func(t *testing.T) {
		sale := &models.Sale{ID: 4, Status: "pending"}
		mockRepo.On("GetByID", ctx, int64(4)).Return(sale, nil).Once()
		err := svc.Returned(ctx, 4)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "somente vendas concluídas podem ser devolvidas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale status canceled", func(t *testing.T) {
		sale := &models.Sale{ID: 5, Status: "canceled"}
		mockRepo.On("GetByID", ctx, int64(5)).Return(sale, nil).Once()
		err := svc.Returned(ctx, 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "somente vendas concluídas podem ser devolvidas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale status returned", func(t *testing.T) {
		sale := &models.Sale{ID: 6, Status: "returned"}
		mockRepo.On("GetByID", ctx, int64(6)).Return(sale, nil).Once()
		err := svc.Returned(ctx, 6)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "somente vendas concluídas podem ser devolvidas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo update", func(t *testing.T) {
		sale := &models.Sale{ID: 7, Status: "completed"}
		mockRepo.On("GetByID", ctx, int64(7)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(errors.New("update error")).Once()

		err := svc.Returned(ctx, 7)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		sale := &models.Sale{ID: 8, Status: "completed"}
		mockRepo.On("GetByID", ctx, int64(8)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(nil).Once()

		err := svc.Returned(ctx, 8)
		assert.NoError(t, err)
		assert.Equal(t, "returned", sale.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_Activate(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()

	t.Run("id inválido", func(t *testing.T) {
		err := svc.Activate(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("sale não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errMsg.ErrNotFound).Once()
		err := svc.Activate(ctx, 1)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo get", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(2)).Return(nil, errors.New("db error")).Once()

		err := svc.Activate(ctx, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale não cancelada nem devolvida", func(t *testing.T) {
		sale := &models.Sale{ID: 3, Status: "active"}
		mockRepo.On("GetByID", ctx, int64(3)).Return(sale, nil).Once()
		err := svc.Activate(ctx, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "somente vendas canceladas ou devolvidas podem ser reativadas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale status completed", func(t *testing.T) {
		sale := &models.Sale{ID: 4, Status: "completed"}
		mockRepo.On("GetByID", ctx, int64(4)).Return(sale, nil).Once()
		err := svc.Activate(ctx, 4)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "somente vendas canceladas ou devolvidas podem ser reativadas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale status pending", func(t *testing.T) {
		sale := &models.Sale{ID: 5, Status: "pending"}
		mockRepo.On("GetByID", ctx, int64(5)).Return(sale, nil).Once()
		err := svc.Activate(ctx, 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "somente vendas canceladas ou devolvidas podem ser reativadas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso com sale cancelada", func(t *testing.T) {
		sale := &models.Sale{ID: 6, Status: "canceled"}
		mockRepo.On("GetByID", ctx, int64(6)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(nil).Once()

		err := svc.Activate(ctx, 6)
		assert.NoError(t, err)
		assert.Equal(t, "active", sale.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso com sale devolvida", func(t *testing.T) {
		sale := &models.Sale{ID: 7, Status: "returned"}
		mockRepo.On("GetByID", ctx, int64(7)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(nil).Once()

		err := svc.Activate(ctx, 7)
		assert.NoError(t, err)
		assert.Equal(t, "active", sale.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro repo update", func(t *testing.T) {
		sale := &models.Sale{ID: 8, Status: "canceled"}
		mockRepo.On("GetByID", ctx, int64(8)).Return(sale, nil).Once()
		mockRepo.On("Update", ctx, sale).Return(errors.New("update error")).Once()

		err := svc.Activate(ctx, 8)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})
}
