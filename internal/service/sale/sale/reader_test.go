package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockSale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleService_GetByID(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)

	t.Run("id inválido", func(t *testing.T) {
		result, err := svc.GetByID(context.Background(), 0)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro repo not found", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(nil, errMsg.ErrNotFound).Once()
		result, err := svc.GetByID(context.Background(), 1)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro genérico do repo", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockRepo.On("GetByID", mock.Anything, int64(3)).Return(nil, expectedErr).Once()
		result, err := svc.GetByID(context.Background(), 3)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		expectedSale := &models.Sale{ID: 2, UserID: utils.Int64Ptr(1)}
		mockRepo.On("GetByID", mock.Anything, int64(2)).Return(expectedSale, nil).Once()
		result, err := svc.GetByID(context.Background(), 2)
		assert.NoError(t, err)
		assert.Equal(t, expectedSale, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByClientID(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)

	t.Run("clientID inválido", func(t *testing.T) {
		result, err := svc.GetByClientID(context.Background(), 0, 10, 0, "sale_date", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("limite inválido", func(t *testing.T) {
		result, err := svc.GetByClientID(context.Background(), 1, 0, 0, "sale_date", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidLimit)
	})

	t.Run("offset inválido", func(t *testing.T) {
		result, err := svc.GetByClientID(context.Background(), 1, 10, -1, "sale_date", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOffset)
	})

	t.Run("ordem inválida", func(t *testing.T) {
		result, err := svc.GetByClientID(context.Background(), 1, 10, 0, "invalid_field", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderField)
	})

	t.Run("direção inválida", func(t *testing.T) {
		result, err := svc.GetByClientID(context.Background(), 1, 10, 0, "sale_date", "invalid")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderDirection)
	})

	t.Run("sucesso", func(t *testing.T) {
		expectedSales := []*models.Sale{
			{ID: 1, UserID: utils.Int64Ptr(1)},
			{ID: 2, UserID: utils.Int64Ptr(1)},
		}
		mockRepo.On("GetByClientID", mock.Anything, int64(1), 10, 0, "sale_date", "asc").
			Return(expectedSales, nil).Once()
		result, err := svc.GetByClientID(context.Background(), 1, 10, 0, "sale_date", "asc")
		assert.NoError(t, err)
		assert.Equal(t, expectedSales, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByUserID(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()

	t.Run("id inválido", func(t *testing.T) {
		result, err := svc.GetByUserID(ctx, 0, 10, 0, "sale_date", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("erro paginação", func(t *testing.T) {
		result, err := svc.GetByUserID(ctx, 1, 0, 0, "sale_date", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidLimit)
	})

	t.Run("erro order field", func(t *testing.T) {
		result, err := svc.GetByUserID(ctx, 1, 10, 0, "invalid_field", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderField)
	})

	t.Run("erro order direction", func(t *testing.T) {
		result, err := svc.GetByUserID(ctx, 1, 10, 0, "sale_date", "invalid")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderDirection)
	})

	t.Run("erro genérico do repo", func(t *testing.T) {
		expectedErr := errors.New("repo error")
		mockRepo.On("GetByUserID", ctx, int64(1), 10, 0, "sale_date", "asc").Return(nil, expectedErr).Once()
		result, err := svc.GetByUserID(ctx, 1, 10, 0, "sale_date", "asc")
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		expectedSales := []*models.Sale{{ID: 1}, {ID: 2}}
		mockRepo.On("GetByUserID", ctx, int64(1), 10, 0, "sale_date", "asc").Return(expectedSales, nil).Once()
		result, err := svc.GetByUserID(ctx, 1, 10, 0, "sale_date", "asc")
		assert.NoError(t, err)
		assert.Equal(t, expectedSales, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_GetByDateRange(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	t.Run("data inválida", func(t *testing.T) {
		result, err := svc.GetByDateRange(ctx, time.Time{}, end, 10, 0, "id", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("start depois do end", func(t *testing.T) {
		result, err := svc.GetByDateRange(ctx, end, start, 10, 0, "id", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidDateRange)
	})

	t.Run("erro paginação", func(t *testing.T) {
		result, err := svc.GetByDateRange(ctx, start, end, 0, 0, "id", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidLimit)
	})

	t.Run("erro order field", func(t *testing.T) {
		result, err := svc.GetByDateRange(ctx, start, end, 10, 0, "invalid", "asc")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderField)
	})

	t.Run("erro order direction", func(t *testing.T) {
		result, err := svc.GetByDateRange(ctx, start, end, 10, 0, "id", "invalid")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOrderDirection)
	})

	t.Run("erro genérico do repo", func(t *testing.T) {
		expectedErr := errors.New("repo error")
		mockRepo.On("GetByDateRange", ctx, start, end, 10, 0, "id", "asc").Return(nil, expectedErr).Once()
		result, err := svc.GetByDateRange(ctx, start, end, 10, 0, "id", "asc")
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		expectedSales := []*models.Sale{{ID: 1}, {ID: 2}}
		mockRepo.On("GetByDateRange", ctx, start, end, 10, 0, "id", "asc").Return(expectedSales, nil).Once()
		result, err := svc.GetByDateRange(ctx, start, end, 10, 0, "id", "asc")
		assert.NoError(t, err)
		assert.Equal(t, expectedSales, result)
		mockRepo.AssertExpectations(t)
	})
}
