package services

import (
	"context"
	"errors"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	user_categories_repositories_mock "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryService_Create(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New()) // logger real

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)

		inputCategory := &models.UserCategory{Name: "NewCategory", Description: "NewDesc"}
		createdCategory := &models.UserCategory{ID: 1, Name: "NewCategory", Description: "NewDesc"}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(cat *models.UserCategory) bool {
			return cat.Name == inputCategory.Name && cat.Description == inputCategory.Description
		})).Return(createdCategory, nil)

		service := NewUserCategoryService(mockRepo, logger)
		category, err := service.Create(context.Background(), inputCategory)

		assert.NoError(t, err)
		assert.Equal(t, createdCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrInvalidCategoryName", func(t *testing.T) {
		mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)
		service := NewUserCategoryService(mockRepo, logger)

		invalidCategory := &models.UserCategory{Name: "   "} // nome só com espaços

		category, err := service.Create(context.Background(), invalidCategory)

		assert.Nil(t, category)
		assert.ErrorIs(t, err, ErrInvalidCategoryName)
	})

	t.Run("ErrorOnCreate", func(t *testing.T) {
		mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)
		inputCategory := &models.UserCategory{Name: "NewCategory", Description: "NewDesc"}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(cat *models.UserCategory) bool {
			return cat.Name == inputCategory.Name && cat.Description == inputCategory.Description
		})).Return(nil, errors.New("db error"))

		service := NewUserCategoryService(mockRepo, logger)
		category, err := service.Create(context.Background(), inputCategory)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar categoria")
		assert.Contains(t, err.Error(), "db error")
		assert.Nil(t, category)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryService_GetAll(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)
		expectedCategories := []*models.UserCategory{
			{ID: 1, Name: "Category1", Description: "Desc1"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedCategories, nil)

		service := NewUserCategoryService(mockRepo, logger)
		categories, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnGetAll", func(t *testing.T) {
		mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)
		mockRepo.On("GetAll", mock.Anything).Return([]*models.UserCategory(nil), errors.New("db error"))

		service := NewUserCategoryService(mockRepo, logger)
		categories, err := service.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "erro ao buscar categorias")
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryService_GetCategoryById(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())
	mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)
	service := NewUserCategoryService(mockRepo, logger)

	t.Run("Success", func(t *testing.T) {
		expectedCategory := &models.UserCategory{ID: 1, Name: "Category1", Description: "Desc1"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedCategory, nil).Once()

		category, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrCategoryIDRequired", func(t *testing.T) {
		category, err := service.GetByID(context.Background(), 0)

		assert.Nil(t, category)
		assert.ErrorIs(t, err, ErrCategoryIDRequired)
	})

	t.Run("NotFound", func(t *testing.T) {
		notFoundErr := errors.New("categoria não encontrada")
		mockRepo.On("GetByID", mock.Anything, int64(2)).Return((*models.UserCategory)(nil), notFoundErr).Once()

		category, err := service.GetByID(context.Background(), 2)

		assert.ErrorIs(t, err, notFoundErr)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ReturnCategoryNotFound", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(4)).Return(nil, ErrCategoryNotFound).Once()

		category, err := service.GetByID(context.Background(), 4)

		assert.ErrorIs(t, err, ErrCategoryNotFound)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ReturnCategoryNotFound_Duplicate", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(4)).Return((*models.UserCategory)(nil), ErrCategoryNotFound).Once()

		category, err := service.GetByID(context.Background(), 4)

		assert.ErrorIs(t, err, ErrCategoryNotFound)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryService_UpdateCategory(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())
	mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)
	service := NewUserCategoryService(mockRepo, logger)

	t.Run("Success", func(t *testing.T) {
		updatedCategory := &models.UserCategory{ID: 1, Name: "UpdatedCategory", Description: "UpdatedDesc"}
		mockRepo.On("Update", mock.Anything, updatedCategory).Return(nil).Once()

		category, err := service.Update(context.Background(), updatedCategory)

		assert.NoError(t, err)
		assert.Equal(t, updatedCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		updatedCategory := &models.UserCategory{ID: 2, Name: "FailCategory", Description: "FailDesc"}
		repoErr := errors.New("erro ao atualizar categoria")
		mockRepo.On("Update", mock.Anything, updatedCategory).Return(repoErr).Once()

		category, err := service.Update(context.Background(), updatedCategory)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar categoria")
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		missingCategory := &models.UserCategory{ID: 4, Name: "NotFound", Description: "NoDesc"}
		mockRepo.On("Update", mock.Anything, missingCategory).Return(user_categories_repositories_mock.ErrCategoryNotFound).Once()

		category, err := service.Update(context.Background(), missingCategory)

		assert.ErrorIs(t, err, user_categories_repositories_mock.ErrCategoryNotFound)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCategoryNil", func(t *testing.T) {
		category, err := service.Update(context.Background(), nil)
		assert.ErrorIs(t, err, ErrInvalidCategory)
		assert.Nil(t, category)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		category := &models.UserCategory{ID: 0}
		result, err := service.Update(context.Background(), category)
		assert.ErrorIs(t, err, ErrCategoryIDRequired)
		assert.Nil(t, result)
	})
}

func TestUserCategoryService_DeleteCategoryById(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())
	mockRepo := new(user_categories_repositories_mock.MockUserCategoryRepository)
	service := NewUserCategoryService(mockRepo, logger)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		repoErr := errors.New("db delete error")
		mockRepo.On("Delete", mock.Anything, int64(2)).Return(repoErr).Once()

		err := service.Delete(context.Background(), 2)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrDeleteCategory))
		assert.ErrorContains(t, err, "db delete error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, ErrCategoryIDRequired)
	})
}
