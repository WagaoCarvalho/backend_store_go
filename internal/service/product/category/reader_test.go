package services

import (
	"context"
	"errors"
	"testing"

	mockProductCat "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryService_GetAll(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)
		expectedCategories := []*models.ProductCategory{
			{ID: 1, Name: "Category1", Description: "Desc1"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedCategories, nil)

		service := NewProductCategoryService(mockRepo)
		categories, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnGetAll", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)
		mockRepo.On("GetAll", mock.Anything).Return([]*models.ProductCategory(nil), errors.New("db error"))

		service := NewProductCategoryService(mockRepo)
		categories, err := service.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "erro ao buscar")
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryService_GetById(t *testing.T) {

	mockRepo := new(mockProductCat.MockProductCategory)
	service := NewProductCategoryService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedCategory := &models.ProductCategory{ID: 1, Name: "Category1", Description: "Desc1"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedCategory, nil).Once()

		category, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrCategoryIDRequired", func(t *testing.T) {
		category, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, category)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("ReturnCategoryNotFound", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(4)).Return(nil, errMsg.ErrNotFound).Once()

		category, err := service.GetByID(context.Background(), 4)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		ctx := context.Background()
		internalErr := errors.New("erro interno do banco")

		mockRepo.On("GetByID", ctx, int64(5)).Return(nil, internalErr).Once()

		category, err := service.GetByID(ctx, 5)

		assert.Nil(t, category)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, "erro interno do banco")
		mockRepo.AssertExpectations(t)
	})

	t.Run("ReturnCategoryNotFound_Duplicate", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(4)).Return((*models.ProductCategory)(nil), errMsg.ErrNotFound).Once()

		category, err := service.GetByID(context.Background(), 4)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})
}
