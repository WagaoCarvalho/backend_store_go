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

func Test_HasProductCategoryRelation(t *testing.T) {

	t.Run("Success_ExistsTrue", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		mockRepo.On("HasProductCategoryRelation", mock.Anything, int64(1), int64(2)).Return(true, nil)

		exists, err := service.HasProductCategoryRelation(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success_ExistsFalse", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		mockRepo.On("HasProductCategoryRelation", mock.Anything, int64(1), int64(3)).Return(false, nil)

		exists, err := service.HasProductCategoryRelation(context.Background(), 1, 3)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidProductID", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		_, err := service.HasProductCategoryRelation(context.Background(), 0, 1)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		_, err := service.HasProductCategoryRelation(context.Background(), 1, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		expectedErr := errors.New("erro no banco de dados")
		mockRepo.On("HasProductCategoryRelation", mock.Anything, int64(1), int64(2)).Return(false, expectedErr)

		exists, err := service.HasProductCategoryRelation(context.Background(), 1, 2)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrRelationCheck)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}
