package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mocksale "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/sale"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	errmsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleService_Create(t *testing.T) {
	t.Run("falha na validação da venda", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		saleModel := &models.Sale{} // inválido: sem UserID, PaymentType etc.

		createdSale, err := service.Create(context.Background(), saleModel)

		assert.Nil(t, createdSale)
		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("sucesso na criação da venda", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		saleModel := &models.Sale{
			UserID:      1,
			SaleDate:    time.Now(),
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		mockRepo.On("Create", mock.Anything, saleModel).Return(saleModel, nil)

		createdSale, err := service.Create(context.Background(), saleModel)

		assert.NoError(t, err)
		assert.Equal(t, saleModel, createdSale)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		saleModel := &models.Sale{
			UserID:      1,
			SaleDate:    time.Now(),
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Create", mock.Anything, saleModel).Return((*models.Sale)(nil), expectedErr)

		createdSale, err := service.Create(context.Background(), saleModel)

		assert.Nil(t, createdSale)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByID(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		result, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errmsg.ErrID)
		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("não encontrado", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return((*models.Sale)(nil), errmsg.ErrNotFound)

		result, err := service.GetByID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errmsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		unexpectedErr := errors.New("erro no banco")
		mockRepo.On("GetByID", mock.Anything, int64(2)).Return((*models.Sale)(nil), unexpectedErr)

		result, err := service.GetByID(context.Background(), 2)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errmsg.ErrGet.Error())
		assert.Contains(t, err.Error(), unexpectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		expected := &models.Sale{ID: 3, UserID: 1, PaymentType: "cash", Status: "active"}
		mockRepo.On("GetByID", mock.Anything, int64(3)).Return(expected, nil)

		result, err := service.GetByID(context.Background(), 3)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByClientID(t *testing.T) {
	t.Run("falha por clientID inválido", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		sales, err := service.GetByClientID(context.Background(), 0, 10, 0, "sale_date", "asc")

		assert.Nil(t, sales)
		assert.Equal(t, err_msg.ErrID, err)
		mockRepo.AssertNotCalled(t, "GetByClientID", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("sucesso ao buscar por clientID", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		clientID := int64(1)
		saleList := []*models.Sale{
			{ID: 1, ClientID: &clientID, UserID: 1, SaleDate: time.Now(), TotalAmount: 100.0},
		}

		mockRepo.On("GetByClientID", mock.Anything, clientID, 10, 0, "sale_date", "asc").Return(saleList, nil)

		result, err := service.GetByClientID(context.Background(), clientID, 10, 0, "sale_date", "asc")

		assert.NoError(t, err)
		assert.Equal(t, saleList, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		clientID := int64(1)
		mockRepo.On("GetByClientID", mock.Anything, clientID, 10, 0, "sale_date", "asc").
			Return(nil, assert.AnError)

		sales, err := service.GetByClientID(context.Background(), clientID, 10, 0, "sale_date", "asc")
		assert.Nil(t, sales)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByUserID(t *testing.T) {
	t.Run("falha por userID inválido", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		sales, err := service.GetByUserID(context.Background(), 0, 10, 0, "sale_date", "asc")

		assert.Nil(t, sales)
		assert.Equal(t, err_msg.ErrID, err)
		mockRepo.AssertNotCalled(t, "GetByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("sucesso ao buscar por userID", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		userID := int64(1)
		saleList := []*models.Sale{
			{ID: 1, ClientID: nil, UserID: userID, SaleDate: time.Now(), TotalAmount: 100.0},
		}

		mockRepo.On("GetByUserID", mock.Anything, userID, 10, 0, "sale_date", "asc").Return(saleList, nil)

		result, err := service.GetByUserID(context.Background(), userID, 10, 0, "sale_date", "asc")

		assert.NoError(t, err)
		assert.Equal(t, saleList, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		userID := int64(1)
		mockRepo.On("GetByUserID", mock.Anything, userID, 10, 0, "sale_date", "asc").
			Return(nil, assert.AnError)

		sales, err := service.GetByUserID(context.Background(), userID, 10, 0, "sale_date", "asc")
		assert.Nil(t, sales)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByStatus(t *testing.T) {
	t.Run("falha por status vazio", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		sales, err := service.GetByStatus(context.Background(), "", 10, 0, "sale_date", "asc")

		assert.Nil(t, sales)
		assert.Equal(t, err_msg.ErrInvalidData, err)
		mockRepo.AssertNotCalled(t, "GetByStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("sucesso ao buscar por status", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		status := "active"
		saleList := []*models.Sale{
			{ID: 1, UserID: 1, SaleDate: time.Now(), TotalAmount: 100.0, Status: status},
		}

		mockRepo.On("GetByStatus", mock.Anything, status, 10, 0, "sale_date", "asc").Return(saleList, nil)

		result, err := service.GetByStatus(context.Background(), status, 10, 0, "sale_date", "asc")

		assert.NoError(t, err)
		assert.Equal(t, saleList, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		status := "active"
		mockRepo.On("GetByStatus", mock.Anything, status, 10, 0, "sale_date", "asc").
			Return(nil, assert.AnError)

		sales, err := service.GetByStatus(context.Background(), status, 10, 0, "sale_date", "asc")
		assert.Nil(t, sales)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByDateRange(t *testing.T) {
	t.Run("falha por datas inválidas", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		sales, err := service.GetByDateRange(context.Background(), time.Time{}, time.Time{}, 10, 0, "sale_date", "asc")

		assert.Nil(t, sales)
		assert.Equal(t, err_msg.ErrInvalidData, err)
		mockRepo.AssertNotCalled(t, "GetByDateRange", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("sucesso ao buscar por range de datas", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		start := time.Now().Add(-24 * time.Hour)
		end := time.Now()
		saleList := []*models.Sale{
			{ID: 1, UserID: 1, SaleDate: start.Add(1 * time.Hour), TotalAmount: 100.0, Status: "active"},
		}

		mockRepo.On("GetByDateRange", mock.Anything, start, end, 10, 0, "sale_date", "asc").Return(saleList, nil)

		result, err := service.GetByDateRange(context.Background(), start, end, 10, 0, "sale_date", "asc")

		assert.NoError(t, err)
		assert.Equal(t, saleList, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		start := time.Now().Add(-24 * time.Hour)
		end := time.Now()
		mockRepo.On("GetByDateRange", mock.Anything, start, end, 10, 0, "sale_date", "asc").
			Return(nil, assert.AnError)

		sales, err := service.GetByDateRange(context.Background(), start, end, 10, 0, "sale_date", "asc")
		assert.Nil(t, sales)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_Update(t *testing.T) {
	t.Run("falha na validação", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		saleModel := &models.Sale{} // inválido

		err := service.Update(context.Background(), saleModel)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("sucesso no update", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		saleModel := &models.Sale{
			ID:          1,
			UserID:      1,
			PaymentType: "cash",
			Status:      "active",
			TotalAmount: 50,
		}

		mockRepo.On("Update", mock.Anything, saleModel).Return(nil)

		err := service.Update(context.Background(), saleModel)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro no repo", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		saleModel := &models.Sale{
			ID:          1,
			UserID:      1,
			PaymentType: "cash",
			Status:      "active",
		}

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Update", mock.Anything, saleModel).Return(expectedErr)

		err := service.Update(context.Background(), saleModel)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), errmsg.ErrUpdate.Error())
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_Delete(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errmsg.ErrID)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro no repo", func(t *testing.T) {
		mockRepo := new(mocksale.MockSaleRepository)
		service := NewSaleService(mockRepo)

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Delete", mock.Anything, int64(2)).Return(expectedErr)

		err := service.Delete(context.Background(), 2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), errmsg.ErrDelete.Error())
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}
