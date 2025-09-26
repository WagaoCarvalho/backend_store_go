package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	repo "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func Test_Create(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
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
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, _, err := service.Create(context.Background(), 0, 1)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)

		_, _, err = service.Create(context.Background(), 1, 0)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
	})

	t.Run("AlreadyExists_ReturnsExisting", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		existing := &models.UserCategoryRelations{UserID: 1, CategoryID: 2}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations{existing}, nil)

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.False(t, wasCreated)
		assert.Equal(t, *existing, *result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_GetByUserIDFails", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return([]*models.UserCategoryRelations{}, errors.New("db error"))

		_, _, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao verificar relação")
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists_ButRelationNotFound", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return([]*models.UserCategoryRelations{
			{UserID: 1, CategoryID: 999},
		}, nil)

		_, _, err := service.Create(context.Background(), 1, 2)

		assert.ErrorIs(t, err, err_msg.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ForeignKeyViolation_ReturnsInvalidForeignKeyError", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		userID := int64(1)
		categoryID := int64(999)

		mockRepo.
			On("Create", mock.Anything, mock.Anything).
			Return(nil, err_msg.ErrInvalidForeignKey)

		rel, wasCreated, err := service.Create(context.Background(), userID, categoryID)

		assert.Nil(t, rel)
		assert.False(t, wasCreated)
		assert.ErrorIs(t, err, err_msg.ErrInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		_, _, err := service.Create(context.Background(), 1, 2)

		assert.ErrorContains(t, err, "erro ao criar")
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetAllRelationsByUserID(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		expected := []*models.UserCategoryRelations{{UserID: 1, CategoryID: 2}}
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := service.GetAllRelationsByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, err := service.GetAllRelationsByUserID(context.Background(), 0)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		expectedErr := errors.New("erro no banco de dados")
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return(([]*models.UserCategoryRelations)(nil), expectedErr)

		result, err := service.GetAllRelationsByUserID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrGet)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_HasUserCategoryRelation(t *testing.T) {

	t.Run("Success_ExistsTrue", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("HasUserCategoryRelation", mock.Anything, int64(1), int64(2)).Return(true, nil)

		exists, err := service.HasUserCategoryRelation(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success_ExistsFalse", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		mockRepo.On("HasUserCategoryRelation", mock.Anything, int64(1), int64(3)).Return(false, nil)

		exists, err := service.HasUserCategoryRelation(context.Background(), 1, 3)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, err := service.HasUserCategoryRelation(context.Background(), 0, 1)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		_, err := service.HasUserCategoryRelation(context.Background(), 1, 0)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(repo.MockUserCategoryRelationRepo)
		service := NewUserCategoryRelationServices(mockRepo)

		expectedErr := errors.New("erro no banco de dados")
		mockRepo.On("HasUserCategoryRelation", mock.Anything, int64(1), int64(2)).Return(false, expectedErr)

		exists, err := service.HasUserCategoryRelation(context.Background(), 1, 2)

		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrRelationCheck)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_Delete(t *testing.T) {

	mockRepo := new(repo.MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := service.Delete(context.Background(), 1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		err := service.Delete(context.Background(), 0, 1)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		err := service.Delete(context.Background(), 1, 0)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
	})

	t.Run("RelationNotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(err_msg.ErrNotFound)

		err := service.Delete(context.Background(), 1, 2)

		assert.ErrorIs(t, err, err_msg.ErrNotFound)
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

	mockRepo := new(repo.MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil)

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		err := service.DeleteAll(context.Background(), 0)

		assert.ErrorIs(t, err, err_msg.ErrZeroID)
	})

	t.Run("DeleteAllError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.DeleteAll(context.Background(), 1)

		assert.ErrorContains(t, err, "erro ao deletar")
		mockRepo.AssertExpectations(t)
	})
}
