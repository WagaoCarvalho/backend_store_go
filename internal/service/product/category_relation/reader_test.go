package services

import (
	"context"
	"errors"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetAllRelationsByProductID(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		expected := []*models.ProductCategoryRelation{{ProductID: 1, CategoryID: 2}}
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidProductID", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		_, err := service.GetAllRelationsByProductID(context.Background(), 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		expectedErr := errors.New("erro no banco de dados")
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(([]*models.ProductCategoryRelation)(nil), expectedErr)

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}
