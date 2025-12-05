package services

import (
	"context"
	"errors"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_GetByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockProduct.ProductMock)

	service := NewProductService(mockRepo)

	t.Run("GetByID com ID inválido", func(t *testing.T) {
		result, err := service.GetByID(ctx, 0)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("sucesso", func(t *testing.T) {

		id := int64(1)
		expectedProduct := &models.Product{
			ID:          id,
			ProductName: "Produto 1",
		}

		mockRepo.
			On("GetByID", ctx, id).
			Return(expectedProduct, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("id inválido", func(t *testing.T) {

		id := int64(0)

		result, err := service.GetByID(ctx, id)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("erro do repositório", func(t *testing.T) {

		id := int64(99)
		mockErr := errors.New("erro ao buscar produto")

		mockRepo.
			On("GetByID", ctx, id).
			Return(nil, mockErr)

		result, err := service.GetByID(ctx, id)

		assert.Nil(t, result)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetVersionByID(t *testing.T) {
	t.Parallel()

	newService := func() (*mockProduct.ProductMock, ProductService) {
		mr := new(mockProduct.ProductMock)

		return mr, NewProductService(mr)
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

	t.Run("ProductNotFound", func(t *testing.T) {
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
