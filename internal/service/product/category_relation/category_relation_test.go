package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func Test_Create(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}
		expected := input

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidIDs", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		// ProductID inválido
		_, err := service.Create(context.Background(), &models.ProductCategoryRelation{ProductID: 0, CategoryID: 1})
		assert.ErrorIs(t, err, errMsg.ErrZeroID)

		// CategoryID inválido
		_, err = service.Create(context.Background(), &models.ProductCategoryRelation{ProductID: 1, CategoryID: 0})
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("NilModel", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		_, err := service.Create(context.Background(), nil)
		assert.ErrorIs(t, err, errMsg.ErrNilModel)
	})

	t.Run("AlreadyExists_ReturnsExisting", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}
		existing := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return([]*models.ProductCategoryRelation{existing}, nil)

		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, *existing, *result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_GetByProductIDFails", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, input).
			Return(nil, errMsg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return([]*models.ProductCategoryRelation(nil), errors.New("db error"))

		_, err := service.Create(context.Background(), input)

		assert.ErrorContains(t, err, "erro ao verificar relação")
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_ButRelationNotFound", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return([]*models.ProductCategoryRelation{
				{ProductID: 1, CategoryID: 999},
			}, nil)

		_, err := service.Create(context.Background(), input)

		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ForeignKeyViolation_ReturnsInvalidForeignKeyError", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 999}

		mockRepo.On("Create", mock.Anything, input).
			Return(nil, errMsg.ErrDBInvalidForeignKey)

		result, err := service.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := NewProductCategoryRelation(mockRepo)

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, input).
			Return(nil, errors.New("db error"))

		_, err := service.Create(context.Background(), input)

		assert.ErrorContains(t, err, "erro ao criar")
		mockRepo.AssertExpectations(t)
	})
}

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

func Test_Delete(t *testing.T) {

	mockRepo := new(mockProduct.MockProductCategoryRelation)
	service := NewProductCategoryRelation(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := service.Delete(context.Background(), 1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidProductID", func(t *testing.T) {
		err := service.Delete(context.Background(), 0, 1)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		err := service.Delete(context.Background(), 1, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("RelationNotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errMsg.ErrNotFound)

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("db error"))

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao deletar")
		mockRepo.AssertExpectations(t)
	})
}

func Test_DeleteAll(t *testing.T) {

	mockRepo := new(mockProduct.MockProductCategoryRelation)
	service := NewProductCategoryRelation(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil)

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidProductID", func(t *testing.T) {
		err := service.DeleteAll(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("DeleteAllError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.DeleteAll(context.Background(), 1)

		assert.ErrorContains(t, err, "erro ao deletar")
		mockRepo.AssertExpectations(t)
	})
}
