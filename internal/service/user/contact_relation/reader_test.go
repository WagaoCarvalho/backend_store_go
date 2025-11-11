package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockContRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserContactRelationServices_GetAllRelationsByUserID(t *testing.T) {
	t.Run("success - retorna lista de relações", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

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
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		result, err := service.GetAllRelationsByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "GetAllRelationsByUserID")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

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
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

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
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

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
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		exists, err := service.HasUserContactRelation(context.Background(), 0, 10)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "HasUserContactRelation")
	})

	t.Run("erro - contactID inválido", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		exists, err := service.HasUserContactRelation(context.Background(), 1, 0)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.EqualError(t, err, err_msg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "HasUserContactRelation")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

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
