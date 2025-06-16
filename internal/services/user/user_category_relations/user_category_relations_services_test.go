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
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		input := &models.UserCategoryRelations{UserID: 1, CategoryID: 2}
		expected := input

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, wasCreated)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidIDs", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, _, err := service.Create(context.Background(), 0, 1)
		assert.ErrorIs(t, err, ErrInvalidUserID)

		_, _, err = service.Create(context.Background(), 1, 0)
		assert.ErrorIs(t, err, ErrInvalidCategoryID)
	})

	t.Run("AlreadyExists_ReturnsExisting", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		existing := &models.UserCategoryRelations{UserID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations{existing}, nil)

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.False(t, wasCreated)
		assert.Equal(t, *existing, *result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_GetByUserIDFails", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return([]*models.UserCategoryRelations{}, errors.New("db error"))

		_, _, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao verificar relação existente")
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_ButRelationNotFound", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, repoMocks.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations{
			{UserID: 1, CategoryID: 999},
		}, nil)

		_, _, err := service.Create(context.Background(), 1, 2)

		assert.ErrorIs(t, err, repoMocks.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		_, _, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao criar relação")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCategoryRelationServices_GetAllRelationsByUserID(t *testing.T) {
	mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expected := []*models.UserCategoryRelations{{UserID: 1, CategoryID: 2}}
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := service.GetAllRelationsByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		_, err := service.GetAllRelationsByUserID(context.Background(), 0)
		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(repoMocks.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		expectedErr := errors.New("erro no banco de dados")
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations(nil), expectedErr)

		result, err := service.GetAllRelationsByUserID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrFetchUserRelations)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
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
