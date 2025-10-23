package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	repo "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func Test_UserContactRelationServices_Create(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}
		expected := input

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, wasCreated)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero UserID", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		result, wasCreated, err := service.Create(context.Background(), 0, 2)

		assert.Error(t, err)
		assert.False(t, wasCreated)
		assert.Nil(t, result)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Zero ContactID", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		result, wasCreated, err := service.Create(context.Background(), 1, 0)

		assert.Error(t, err)
		assert.False(t, wasCreated)
		assert.Nil(t, result)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Duplicate Relation Exists", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}
		existing := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return([]*models.UserContactRelation{existing}, nil)

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.False(t, wasCreated)
		assert.Equal(t, existing, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Duplicate Relation Exists but not found in GetAll", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return([]*models.UserContactRelation{}, nil)

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.Error(t, err)
		assert.False(t, wasCreated)
		assert.Nil(t, result)
		assert.EqualError(t, err, err_msg.ErrRelationExists.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Duplicate Relation Exists with repo GetAll error", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return(nil, errors.New("db error"))

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.Error(t, err)
		assert.False(t, wasCreated)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), err_msg.ErrRelationCheck.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Foreign Key", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, err_msg.ErrDBInvalidForeignKey)

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.Error(t, err)
		assert.False(t, wasCreated)
		assert.Nil(t, result)
		assert.EqualError(t, err, err_msg.ErrDBInvalidForeignKey.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unexpected Repo Error", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errors.New("unexpected"))

		result, wasCreated, err := service.Create(context.Background(), 1, 2)

		assert.Error(t, err)
		assert.False(t, wasCreated)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), err_msg.ErrCreate.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserContactRelationServices_GetAllRelationsByUserID(t *testing.T) {
	t.Run("success - retorna lista de relações", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		expectedRelations := []*models.UserContactRelation{
			{UserID: userID, ContactID: 10, CreatedAt: time.Now()},
			{UserID: userID, ContactID: 20, CreatedAt: time.Now()},
		}

		mockRepo.On("GetAllRelationsByUserID", mock.Anything, userID).
			Return(expectedRelations, nil)

		result, err := service.GetAllRelationsByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedRelations, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - userID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		result, err := service.GetAllRelationsByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "GetAllRelationsByUserID")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(2)
		expectedErr := errors.New("db error")

		mockRepo.On("GetAllRelationsByUserID", mock.Anything, userID).
			Return(nil, expectedErr)

		result, err := service.GetAllRelationsByUserID(context.Background(), userID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, err_msg.ErrGet)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserContactRelationServices_HasUserContactRelation(t *testing.T) {
	t.Run("success - relação existe", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		contactID := int64(2)

		mockRepo.On("HasUserContactRelation", mock.Anything, userID, contactID).
			Return(true, nil)

		exists, err := service.HasUserContactRelation(context.Background(), userID, contactID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - relação não existe", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		contactID := int64(99)

		mockRepo.On("HasUserContactRelation", mock.Anything, userID, contactID).
			Return(false, nil)

		exists, err := service.HasUserContactRelation(context.Background(), userID, contactID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - userID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		exists, err := service.HasUserContactRelation(context.Background(), 0, 10)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "HasUserContactRelation")
	})

	t.Run("erro - contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		exists, err := service.HasUserContactRelation(context.Background(), 1, 0)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "HasUserContactRelation")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		contactID := int64(2)
		expectedErr := errors.New("db error")

		mockRepo.On("HasUserContactRelation", mock.Anything, userID, contactID).
			Return(false, expectedErr)

		exists, err := service.HasUserContactRelation(context.Background(), userID, contactID)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrRelationCheck)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserContactRelationServices_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		contactID := int64(2)

		mockRepo.On("Delete", mock.Anything, userID, contactID).
			Return(nil)

		err := service.Delete(context.Background(), userID, contactID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - userID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		err := service.Delete(context.Background(), 0, 2)

		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro - contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		err := service.Delete(context.Background(), 1, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro - relação não encontrada", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		contactID := int64(2)

		mockRepo.On("Delete", mock.Anything, userID, contactID).
			Return(err_msg.ErrNotFound)

		err := service.Delete(context.Background(), userID, contactID)

		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado do repositório", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		contactID := int64(2)
		expectedErr := errors.New("db error")

		mockRepo.On("Delete", mock.Anything, userID, contactID).
			Return(expectedErr)

		err := service.Delete(context.Background(), userID, contactID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, err_msg.ErrDelete)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserContactRelationServices_DeleteAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)

		mockRepo.On("DeleteAll", mock.Anything, userID).
			Return(nil)

		err := service.DeleteAll(context.Background(), userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - userID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		err := service.DeleteAll(context.Background(), 0)

		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "DeleteAll")
	})

	t.Run("erro inesperado do repositório", func(t *testing.T) {
		mockRepo := new(repo.MockUserContactRelationRepo)
		service := NewUserContactRelation(mockRepo)

		userID := int64(1)
		expectedErr := errors.New("db error")

		mockRepo.On("DeleteAll", mock.Anything, userID).
			Return(expectedErr)

		err := service.DeleteAll(context.Background(), userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, err_msg.ErrDelete)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}
