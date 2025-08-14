package services

import (
	"context"
	"errors"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_categories"
	repo_mock "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryService_Create(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo_mock.MockUserCategoryRepository)

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
		mockRepo := new(repo_mock.MockUserCategoryRepository)
		service := NewUserCategoryService(mockRepo, logger)

		invalidCategory := &models.UserCategory{Name: "   "} // nome só com espaços

		category, err := service.Create(context.Background(), invalidCategory)

		assert.Nil(t, category)
		assert.ErrorIs(t, err, ErrInvalidCategoryName)
	})

	t.Run("ErrorOnCreate", func(t *testing.T) {
		mockRepo := new(repo_mock.MockUserCategoryRepository)
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
		mockRepo := new(repo_mock.MockUserCategoryRepository)
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
		mockRepo := new(repo_mock.MockUserCategoryRepository)
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

func TestUserCategoryService_GetById(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())
	mockRepo := new(repo_mock.MockUserCategoryRepository)
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

	t.Run("ReturnCategoryNotFound", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(4)).Return(nil, ErrCategoryNotFound).Once()

		category, err := service.GetByID(context.Background(), 4)

		assert.ErrorIs(t, err, ErrCategoryNotFound)
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
		assert.ErrorIs(t, err, ErrFetchCategory)
		assert.ErrorContains(t, err, "erro interno do banco")
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

func TestUserCategoryService_Update(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())
	mockRepo := new(repo_mock.MockUserCategoryRepository)
	service := NewUserCategoryService(mockRepo, logger)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		updatedCategory := &models.UserCategory{ID: 1, Name: "UpdatedCategory", Description: "UpdatedDesc"}

		mockRepo.On("GetByID", ctx, int64(1)).Return(updatedCategory, nil).Once()
		mockRepo.On("Update", ctx, updatedCategory).Return(nil).Once()

		category, err := service.Update(ctx, updatedCategory)

		assert.NoError(t, err)
		assert.Equal(t, updatedCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorWhileFetchingBeforeUpdate", func(t *testing.T) {
		ctx := context.Background()
		category := &models.UserCategory{
			ID:          7,
			Name:        "ErroDB",
			Description: "Erro simulado no GetByID",
		}

		dbErr := errors.New("erro no banco")

		mockRepo.On("GetByID", ctx, int64(7)).Return(nil, dbErr).Once()

		result, err := service.Update(ctx, category)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCheckBeforeUpdate)
		assert.ErrorContains(t, err, "erro no banco")
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		ctx := context.Background()
		updatedCategory := &models.UserCategory{ID: 2, Name: "FailCategory", Description: "FailDesc"}
		repoErr := errors.New("erro ao atualizar categoria")

		mockRepo.On("GetByID", ctx, int64(2)).Return(updatedCategory, nil).Once()
		mockRepo.On("Update", ctx, updatedCategory).Return(repoErr).Once()

		category, err := service.Update(ctx, updatedCategory)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar categoria")
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		ctx := context.Background()
		invalidCategory := &models.UserCategory{
			ID:          9,
			Name:        "", // Deve causar falha de validação
			Description: "Sem nome",
		}

		result, err := service.Update(ctx, invalidCategory)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "campo obrigatório") // Corrigido
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		ctx := context.Background()
		missingCategory := &models.UserCategory{ID: 4, Name: "NotFound", Description: "NoDesc"}

		mockRepo.On("GetByID", ctx, int64(4)).Return(nil, repo_mock.ErrCategoryNotFound).Once()

		category, err := service.Update(ctx, missingCategory)

		assert.ErrorIs(t, err, ErrCategoryNotFound) // usar o erro do serviço
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		category := &models.UserCategory{ID: 0}
		result, err := service.Update(context.Background(), category)
		assert.ErrorIs(t, err, ErrCategoryIDRequired)
		assert.Nil(t, result)
	})
}

func TestUserCategoryService_Delete(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())
	mockRepo := new(repo_mock.MockUserCategoryRepository)
	service := NewUserCategoryService(mockRepo, logger)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		category := &models.UserCategory{ID: 1, Name: "Categoria", Description: "Desc"}

		mockRepo.On("GetByID", ctx, int64(1)).Return(category, nil).Once()
		mockRepo.On("Delete", ctx, int64(1)).Return(nil).Once()

		err := service.Delete(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID repository error", func(t *testing.T) {
		ctx := context.Background()
		id := int64(10)
		dbErr := errors.New("erro inesperado no banco de dados")

		mockRepo.On("GetByID", ctx, id).Return(nil, dbErr).Once()

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.ErrorContains(t, err, ErrFetchCategory.Error())
		assert.ErrorContains(t, err, dbErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, int64(3)).Return(nil, repo.ErrCategoryNotFound).Once()

		err := service.Delete(ctx, 3)

		assert.ErrorIs(t, err, ErrCategoryNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		ctx := context.Background()
		category := &models.UserCategory{ID: 2, Name: "Categoria", Description: "Desc"}
		repoErr := errors.New("db delete error")

		mockRepo.On("GetByID", ctx, int64(2)).Return(category, nil).Once()
		mockRepo.On("Delete", ctx, int64(2)).Return(repoErr).Once()

		err := service.Delete(ctx, 2)

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
