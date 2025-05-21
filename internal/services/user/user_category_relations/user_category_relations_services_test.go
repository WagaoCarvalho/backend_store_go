package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repoMocks "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
)

func TestUserCategoryRelationServices_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		input := models.UserCategoryRelations{UserID: 1, CategoryID: 2}
		expected := input

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

		result, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, expected, *result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidIDs", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, err := service.Create(context.Background(), 0, 1)
		assert.ErrorIs(t, err, ErrInvalidUserID)

		_, err = service.Create(context.Background(), 1, 0)
		assert.ErrorIs(t, err, ErrInvalidCategoryID)
	})

	t.Run("AlreadyExists_ReturnsExisting", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		existing := models.UserCategoryRelations{UserID: 1, CategoryID: 2}
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]models.UserCategoryRelations{existing}, nil)

		result, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, existing, *result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_GetByUserIDFails", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]models.UserCategoryRelations(nil), errors.New("db error"))

		_, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao verificar relação existente")
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_ButRelationNotFound", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]models.UserCategoryRelations{
			{UserID: 1, CategoryID: 999},
		}, nil)

		_, err := service.Create(context.Background(), 1, 2)

		assert.ErrorIs(t, err, repoMocks.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		_, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao criar relação")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAll_Success(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	expected := []models.UserCategoryRelations{{UserID: 1, CategoryID: 2}}
	mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(expected, nil)

	result, err := service.GetAll(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAll_InvalidUserID(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	_, err := service.GetAll(context.Background(), 0)
	assert.ErrorIs(t, err, ErrInvalidUserID)
}

func TestUserCategoryRelationServices_GetRelations(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		expected := []models.UserCategoryRelations{{UserID: 1, CategoryID: 2}}
		mockRepo.On("GetByCategoryID", mock.Anything, int64(2)).Return(expected, nil)

		result, err := service.GetRelations(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, err := service.GetRelations(context.Background(), 0)
		assert.ErrorIs(t, err, ErrInvalidCategoryID)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("GetByCategoryID", mock.Anything, int64(99)).Return([]models.UserCategoryRelations(nil), errors.New("db error"))

		_, err := service.GetRelations(context.Background(), 99)
		assert.ErrorContains(t, err, "erro ao buscar relações da categoria")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryRelationServices_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := service.Delete(context.Background(), 1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

		err := service.Delete(context.Background(), 0, 1)
		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

		err := service.Delete(context.Background(), 1, 0)
		assert.ErrorIs(t, err, ErrInvalidCategoryID)
	})

	t.Run("RelationNotFound", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(repoMocks.ErrRelationNotFound)

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorIs(t, err, repoMocks.ErrRelationNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteError", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("db error"))

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao deletar relação")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryRelationServices_DeleteAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil)

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

		err := service.DeleteAll(context.Background(), 0)

		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("DeleteAllError", func(t *testing.T) {
		mockRepo := new(MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.DeleteAll(context.Background(), 1)

		assert.ErrorContains(t, err, "erro ao deletar todas as relações do usuário")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByCategoryID_SuccessTrue(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	userID := int64(1)
	categoryID := int64(10)

	relations := []models.UserCategoryRelations{
		{UserID: userID, CategoryID: 5},
		{UserID: userID, CategoryID: categoryID}, // match
	}

	mockRepo.On("GetByUserID", mock.Anything, userID).Return(relations, nil)

	result, err := service.GetByCategoryID(context.Background(), userID, categoryID)

	assert.NoError(t, err)
	assert.True(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByCategoryID_SuccessFalse(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	userID := int64(1)
	categoryID := int64(99)

	relations := []models.UserCategoryRelations{
		{UserID: userID, CategoryID: 1},
		{UserID: userID, CategoryID: 2},
	}

	mockRepo.On("GetByUserID", mock.Anything, userID).Return(relations, nil)

	result, err := service.GetByCategoryID(context.Background(), userID, categoryID)

	assert.NoError(t, err)
	assert.False(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByCategoryID_InvalidUserID(t *testing.T) {
	service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

	result, err := service.GetByCategoryID(context.Background(), 0, 1)

	assert.ErrorIs(t, err, ErrInvalidUserID)
	assert.False(t, result)
}

func TestGetByCategoryID_InvalidCategoryID(t *testing.T) {
	service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

	result, err := service.GetByCategoryID(context.Background(), 1, 0)

	assert.ErrorIs(t, err, ErrInvalidCategoryID)
	assert.False(t, result)
}

func TestGetByCategoryID_RepoError(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	mockRepo.On("GetByUserID", mock.Anything, int64(1)).
		Return([]models.UserCategoryRelations(nil), errors.New("erro inesperado"))

	result, err := service.GetByCategoryID(context.Background(), 1, 2)

	assert.Error(t, err)
	assert.False(t, result)
	mockRepo.AssertExpectations(t)
}
