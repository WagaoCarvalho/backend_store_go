package services

import (
	"context"
	"errors"
	"testing"

	mockRepo "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryRelationService_GetAllRelationsByProductID(t *testing.T) {
	setup := func() (*mockRepo.MockProductCategoryRelation, *productCategoryRelationService) {
		mockRepo := new(mockRepo.MockProductCategoryRelation)
		service := &productCategoryRelationService{repo: mockRepo}
		return mockRepo, service
	}

	t.Run("Success - cria nova instância do serviço", func(t *testing.T) {
		mockRepo := new(mockRepo.MockProductCategoryRelation)

		service := NewProductCategoryRelation(mockRepo)

		assert.NotNil(t, service)

		// Verifica se o tipo retornado é o esperado
		impl, ok := service.(*productCategoryRelationService)
		assert.True(t, ok, "Expected service to be of type *productCategoryRelationService")
		assert.NotNil(t, impl)
		assert.Equal(t, mockRepo, impl.repo)
	})

	t.Run("Success - retorna relações", func(t *testing.T) {
		mockRepo, service := setup()

		expected := []*models.ProductCategoryRelation{
			{ProductID: 1, CategoryID: 2},
			{ProductID: 1, CategoryID: 3},
		}
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(expected, nil).Once()

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - retorna slice vazio quando não há relações", func(t *testing.T) {
		mockRepo, service := setup()

		emptySlice := []*models.ProductCategoryRelation{}
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(emptySlice, nil).Once()

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Equal(t, emptySlice, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - ID inválido (zero)", func(t *testing.T) {
		mockRepo, service := setup()

		result, err := service.GetAllRelationsByProductID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetAllRelationsByProductID")
	})

	t.Run("Error - ID inválido (negativo)", func(t *testing.T) {
		mockRepo, service := setup()

		result, err := service.GetAllRelationsByProductID(context.Background(), -5)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "GetAllRelationsByProductID")
	})

	t.Run("Error - repositório retorna ErrNotFound", func(t *testing.T) {
		mockRepo, service := setup()

		// IMPORTANTE: Retornar nil explicitamente como []*models.ProductCategoryRelation
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(999)).Return(
			([]*models.ProductCategoryRelation)(nil),
			errMsg.ErrNotFound,
		).Once()

		result, err := service.GetAllRelationsByProductID(context.Background(), 999)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		// Verifica que NÃO foi encapsulado
		assert.NotErrorIs(t, err, errMsg.ErrGet)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repositório retorna erro genérico", func(t *testing.T) {
		mockRepo, service := setup()

		expectedErr := errors.New("erro de conexão com o banco")
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(
			([]*models.ProductCategoryRelation)(nil),
			expectedErr,
		).Once()

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.Nil(t, result)
		// Propaga o erro original, não encapsula
		assert.Equal(t, expectedErr, err)
		assert.ErrorContains(t, err, "erro de conexão")
		assert.NotErrorIs(t, err, errMsg.ErrGet)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repositório retorna ErrDBInvalidForeignKey", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(
			([]*models.ProductCategoryRelation)(nil),
			errMsg.ErrDBInvalidForeignKey,
		).Once()

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrGet)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - repo retorna nil, service converte para slice vazio", func(t *testing.T) {
		mockRepo, service := setup()

		// Mock retorna nil explicitamente como []*models.ProductCategoryRelation
		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(
			([]*models.ProductCategoryRelation)(nil),
			nil,
		).Once()

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Equal(t, []*models.ProductCategoryRelation{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repositório retorna ErrRelationExists", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetAllRelationsByProductID", mock.Anything, int64(1)).Return(
			([]*models.ProductCategoryRelation)(nil),
			errMsg.ErrRelationExists,
		).Once()

		result, err := service.GetAllRelationsByProductID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrGet)
		mockRepo.AssertExpectations(t)
	})
}
