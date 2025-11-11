package services

import (
	"context"
	"errors"
	"testing"

	mockSupplierCategory "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierCategoryService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockSupplierCategory.MockSupplierCategory)
		service := NewSupplierCategory(mockRepo)

		expected := &models.SupplierCategory{ID: 1, Name: "Eletrônicos"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := service.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(mockSupplierCategory.MockSupplierCategory)
		service := NewSupplierCategory(mockRepo)

		result, err := service.GetByID(ctx, -1)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mockSupplierCategory.MockSupplierCategory)
		service := NewSupplierCategory(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(999)).Return((*models.SupplierCategory)(nil), errMsg.ErrNotFound)

		result, err := service.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockSupplierCategory.MockSupplierCategory)
		service := NewSupplierCategory(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(2)).Return((*models.SupplierCategory)(nil), errors.New("erro no banco"))

		result, err := service.GetByID(ctx, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockSupplierCategory.MockSupplierCategory)
		service := NewSupplierCategory(mockRepo)

		expectedCategories := []*models.SupplierCategory{
			{ID: 1, Name: "Categoria A"},
			{ID: 2, Name: "Categoria B"},
		}

		mockRepo.On("GetAll", ctx).Return(expectedCategories, nil)

		categories, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockSupplierCategory.MockSupplierCategory)
		service := NewSupplierCategory(mockRepo)

		mockRepo.On("GetAll", ctx).Return(([]*models.SupplierCategory)(nil), errors.New("erro ao buscar categorias"))

		categories, err := service.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "erro ao buscar categorias")
		mockRepo.AssertExpectations(t)
	})
}
