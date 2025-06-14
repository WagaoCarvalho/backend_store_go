package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repoMocks "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
)

func TestUserCategoryRelationServices_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		input := &models.UserCategoryRelations{UserID: 1, CategoryID: 2}
		expected := input

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

		result, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidIDs", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, err := service.Create(context.Background(), 0, 1)
		assert.ErrorIs(t, err, ErrInvalidUserID)

		_, err = service.Create(context.Background(), 1, 0)
		assert.ErrorIs(t, err, ErrInvalidCategoryID)
	})

	t.Run("AlreadyExists_ReturnsExisting", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		existing := &models.UserCategoryRelations{UserID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations{existing}, nil) // slice de ponteiros

		result, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, *existing, *result) // desreferenciando para comparar struct
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_GetByUserIDFails", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations(nil), errors.New("db error"))

		_, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao verificar relação existente")
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_ButRelationNotFound", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations{
			{UserID: 1, CategoryID: 999},
		}, nil)

		_, err := service.Create(context.Background(), 1, 2)

		assert.ErrorIs(t, err, repoMocks.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		_, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao criar relação")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryRelationServices_GetAll(t *testing.T) {
	mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expected := []*models.UserCategoryRelations{{UserID: 1, CategoryID: 2}}
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := service.GetAll(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		_, err := service.GetAll(context.Background(), 0)
		assert.ErrorIs(t, err, ErrInvalidUserID)
	})
}

func TestUserCategoryRelationServices_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := service.Delete(context.Background(), 1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		service := NewUserCategoryRelationServices(new(repoMocks.MockUserCategoryRelationRepo))

		err := service.Delete(context.Background(), 0, 1)
		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		service := NewUserCategoryRelationServices(new(repoMocks.MockUserCategoryRelationRepo))

		err := service.Delete(context.Background(), 1, 0)
		assert.ErrorIs(t, err, ErrInvalidCategoryID)
	})

	t.Run("RelationNotFound", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(repoMocks.ErrRelationNotFound)

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorIs(t, err, repoMocks.ErrRelationNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteError", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("db error"))

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao deletar relação")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryRelationServices_DeleteAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil)

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		service := NewUserCategoryRelationServices(new(repoMocks.MockUserCategoryRelationRepo))

		err := service.DeleteAll(context.Background(), 0)

		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("DeleteAllError", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.DeleteAll(context.Background(), 1)

		assert.ErrorContains(t, err, "erro ao deletar todas as relações do usuário")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByCategoryID(t *testing.T) {
	mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	t.Run("Success_True", func(t *testing.T) {
		userID := int64(1)
		categoryID := int64(10)

		relations := []*models.UserCategoryRelations{
			{UserID: userID, CategoryID: 5},
			{UserID: userID, CategoryID: categoryID}, // match
		}

		mockRepo.On("GetByUserID", mock.Anything, userID).Return(relations, nil)

		result, err := service.GetByCategoryID(context.Background(), userID, categoryID)

		assert.NoError(t, err)
		assert.True(t, result)
		mockRepo.AssertExpectations(t)
		mockRepo.ExpectedCalls = nil // limpa chamadas anteriores
	})

	t.Run("Success_False", func(t *testing.T) {
		userID := int64(1)
		categoryID := int64(99)

		relations := []*models.UserCategoryRelations{
			{UserID: userID, CategoryID: 1},
			{UserID: userID, CategoryID: 2},
		}

		mockRepo.On("GetByUserID", mock.Anything, userID).Return(relations, nil)

		result, err := service.GetByCategoryID(context.Background(), userID, categoryID)

		assert.NoError(t, err)
		assert.False(t, result)
		mockRepo.AssertExpectations(t)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		result, err := service.GetByCategoryID(context.Background(), 0, 1)

		assert.ErrorIs(t, err, ErrInvalidUserID)
		assert.False(t, result)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		result, err := service.GetByCategoryID(context.Background(), 1, 0)

		assert.ErrorIs(t, err, ErrInvalidCategoryID)
		assert.False(t, result)
	})

	t.Run("RepoError", func(t *testing.T) {
		userID := int64(1)

		mockRepo.On("GetByUserID", mock.Anything, userID).
			Return([]*models.UserCategoryRelations(nil), errors.New("erro inesperado"))

		result, err := service.GetByCategoryID(context.Background(), userID, 2)

		assert.Error(t, err)
		assert.False(t, result)
		mockRepo.AssertExpectations(t)
		mockRepo.ExpectedCalls = nil
	})
}

func TestUserCategoryRelationService_GetByUserID(t *testing.T) {
	t.Run("success - retorna relações do usuário", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		expected := []*models.UserCategoryRelations{
			{UserID: 1, CategoryID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{UserID: 1, CategoryID: 11, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		mockRepo.
			On("GetByUserID", mock.Anything, int64(1)).
			Return(expected, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - userID inválido", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidUserID, err)
	})

	t.Run("error - falha no repositório", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.
			On("GetByUserID", mock.Anything, int64(99)).
			Return([]*models.UserCategoryRelations(nil), errors.New("db error"))

		result, err := service.GetByUserID(context.Background(), 99)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}
