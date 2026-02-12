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

func TestProductCategoryRelationService_Create(t *testing.T) {
	setup := func() (*mockProduct.MockProductCategoryRelation, *productCategoryRelationService) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := &productCategoryRelationService{repo: mockRepo}
		return mockRepo, service
	}

	t.Run("Success - cria relação com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}
		expected := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil).Once()

		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - modelo nil", func(t *testing.T) {
		mockRepo, service := setup()

		result, err := service.Create(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNilModel)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - ProductID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.ProductCategoryRelation{ProductID: 0, CategoryID: 1}
		result, err := service.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - CategoryID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 0}
		result, err := service.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - relação já existe", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrRelationExists).Once()

		result, err := service.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - chave estrangeira inválida", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 999}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		result, err := service.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - erro genérico do repositório", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.ProductCategoryRelation{ProductID: 1, CategoryID: 2}
		repoErr := errors.New("erro de conexão")

		mockRepo.On("Create", mock.Anything, input).Return(nil, repoErr).Once()

		result, err := service.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, repoErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryRelationService_Delete(t *testing.T) {
	setup := func() (*mockProduct.MockProductCategoryRelation, *productCategoryRelationService) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := &productCategoryRelationService{repo: mockRepo}
		return mockRepo, service
	}

	t.Run("Success - deleta relação com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil).Once()

		err := service.Delete(context.Background(), 1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - ProductID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		err := service.Delete(context.Background(), 0, 1)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Error - CategoryID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		err := service.Delete(context.Background(), 1, 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Error - relação não encontrada", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errMsg.ErrNotFound).Once()

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - erro genérico do repositório", func(t *testing.T) {
		mockRepo, service := setup()

		repoErr := errors.New("erro de conexão")
		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(repoErr).Once()

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, repoErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryRelationService_DeleteAll(t *testing.T) {
	setup := func() (*mockProduct.MockProductCategoryRelation, *productCategoryRelationService) {
		mockRepo := new(mockProduct.MockProductCategoryRelation)
		service := &productCategoryRelationService{repo: mockRepo}
		return mockRepo, service
	}

	t.Run("Success - deleta todas as relações", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil).Once()

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - produto sem relações (não é erro)", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil).Once()

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - ProductID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		err := service.DeleteAll(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "DeleteAll")
	})

	t.Run("Error - erro genérico do repositório", func(t *testing.T) {
		mockRepo, service := setup()

		repoErr := errors.New("erro de conexão")
		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(repoErr).Once()

		err := service.DeleteAll(context.Background(), 1)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, repoErr.Error())
		mockRepo.AssertExpectations(t)
	})
}
