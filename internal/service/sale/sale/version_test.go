package services

import (
	"context"
	"errors"
	"testing"

	mockSale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleService_GetVersionByID(t *testing.T) {
	t.Parallel()

	newService := func() (*mockSale.MockSale, SaleService) {
		mr := new(mockSale.MockSale)

		return mr, NewSaleService(mr)
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := newService()

		invalidID := int64(0)
		version, err := service.GetVersionByID(context.Background(), invalidID)

		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetVersionByID")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockRepo, service := newService()
		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil)

		version, err := service.GetVersionByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), version)

		mockRepo.AssertExpectations(t)
	})

	t.Run("SaleNotFound", func(t *testing.T) {
		t.Parallel()

		mockRepo, service := newService()
		mockRepo.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), errMsg.ErrNotFound)

		version, err := service.GetVersionByID(context.Background(), 2)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)

		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		t.Parallel()

		mockRepo, service := newService()
		mockRepo.On("GetVersionByID", mock.Anything, int64(3)).Return(int64(0), errors.New("db failure"))

		version, err := service.GetVersionByID(context.Background(), 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db failure")
		assert.True(t, errors.Is(err, errMsg.ErrVersionConflict))
		assert.Equal(t, int64(0), version)

		mockRepo.AssertExpectations(t)
	})
}
