package services_test

import (
	"context"
	"errors"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(services.MockUserCategoryRepository)
		inputCategory := models.UserCategory{Name: "NewCategory", Description: "NewDesc"}
		createdCategory := models.UserCategory{ID: 1, Name: "NewCategory", Description: "NewDesc"}

		mockRepo.On("Create", mock.Anything, inputCategory).Return(createdCategory, nil)

		service := services.NewUserCategoryService(mockRepo)
		category, err := service.Create(context.Background(), inputCategory)

		assert.NoError(t, err)
		assert.Equal(t, createdCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnCreate", func(t *testing.T) {
		mockRepo := new(services.MockUserCategoryRepository)
		inputCategory := models.UserCategory{Name: "NewCategory", Description: "NewDesc"}

		mockRepo.On("Create", mock.Anything, inputCategory).Return(models.UserCategory{}, errors.New("db error"))

		service := services.NewUserCategoryService(mockRepo)
		category, err := service.Create(context.Background(), inputCategory)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar categoria")
		assert.Contains(t, err.Error(), "db error")
		assert.Equal(t, models.UserCategory{}, category)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryService_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(services.MockUserCategoryRepository)
		expectedCategories := []models.UserCategory{
			{ID: 1, Name: "Category1", Description: "Desc1"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedCategories, nil)

		service := services.NewUserCategoryService(mockRepo)
		categories, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnGetAll", func(t *testing.T) {
		mockRepo := new(services.MockUserCategoryRepository)
		mockRepo.On("GetAll", mock.Anything).Return([]models.UserCategory(nil), errors.New("db error"))

		service := services.NewUserCategoryService(mockRepo)
		categories, err := service.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "erro ao buscar categorias")
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryService_GetCategoryById(t *testing.T) {
	mockRepo := new(services.MockUserCategoryRepository)
	service := services.NewUserCategoryService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedCategory := models.UserCategory{ID: 1, Name: "Category1", Description: "Desc1"}
		mockRepo.On("GetById", mock.Anything, int64(1)).Return(expectedCategory, nil).Once()

		category, err := service.GetById(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		notFoundErr := errors.New("categoria não encontrada")
		mockRepo.On("GetById", mock.Anything, int64(2)).Return(models.UserCategory{}, notFoundErr).Once()

		category, err := service.GetById(context.Background(), 2)

		assert.ErrorIs(t, err, notFoundErr)
		assert.Equal(t, models.UserCategory{}, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ReturnCategoryNotFound", func(t *testing.T) {
		mockRepo := new(services.MockUserCategoryRepository)
		service := services.NewUserCategoryService(mockRepo)

		mockRepo.On("GetById", mock.Anything, int64(4)).Return(models.UserCategory{}, services.ErrCategoryNotFound).Once()

		category, err := service.GetById(context.Background(), 4)

		assert.ErrorIs(t, err, services.ErrCategoryNotFound)
		assert.Equal(t, models.UserCategory{}, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GenericError", func(t *testing.T) {
		genericErr := errors.New("db error")
		mockRepo.On("GetById", mock.Anything, int64(3)).Return(models.UserCategory{}, genericErr).Once()

		category, err := service.GetById(context.Background(), 3)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao buscar categoria")
		assert.Equal(t, models.UserCategory{}, category)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryService_UpdateCategory(t *testing.T) {
	mockRepo := new(services.MockUserCategoryRepository)
	service := services.NewUserCategoryService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		updatedCategory := models.UserCategory{ID: 1, Name: "UpdatedCategory", Description: "UpdatedDesc", Version: 1}
		mockRepo.On("Update", mock.Anything, updatedCategory).Return(updatedCategory, nil).Once()

		category, err := service.Update(context.Background(), updatedCategory)

		assert.NoError(t, err)
		assert.Equal(t, updatedCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		updatedCategory := models.UserCategory{ID: 2, Name: "FailCategory", Description: "FailDesc", Version: 1}
		repoErr := errors.New("db update error")
		mockRepo.On("Update", mock.Anything, updatedCategory).Return(models.UserCategory{}, repoErr).Once()

		category, err := service.Update(context.Background(), updatedCategory)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar categoria")
		assert.Equal(t, models.UserCategory{}, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("VersionConflict", func(t *testing.T) {
		conflictCategory := models.UserCategory{ID: 3, Name: "Conflict", Description: "ConflictDesc", Version: 2}
		mockRepo.On("Update", mock.Anything, conflictCategory).Return(models.UserCategory{}, repositories.ErrVersionConflict).Once()

		category, err := service.Update(context.Background(), conflictCategory)

		assert.ErrorIs(t, err, repositories.ErrVersionConflict)
		assert.Equal(t, models.UserCategory{}, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		missingCategory := models.UserCategory{ID: 4, Name: "NotFound", Description: "NoDesc", Version: 1}
		mockRepo.On("Update", mock.Anything, missingCategory).Return(models.UserCategory{}, repositories.ErrCategoryNotFound).Once()

		category, err := service.Update(context.Background(), missingCategory)

		assert.ErrorIs(t, err, repositories.ErrCategoryNotFound)
		assert.Equal(t, models.UserCategory{}, category)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryService_DeleteCategoryById(t *testing.T) {
	mockRepo := new(services.MockUserCategoryRepository)
	service := services.NewUserCategoryService(mockRepo)

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
		assert.ErrorContains(t, err, "erro ao deletar categoria")
		mockRepo.AssertExpectations(t)
	})
}
