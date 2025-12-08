package services

import (
	"context"
	"errors"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
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
