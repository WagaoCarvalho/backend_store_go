package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockContRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func TestUserContactRelation_Create(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}
		expected := *input

		mockRepo.On("Create", mock.Anything, input).Return(&expected, nil)

		result, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NilModel", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		result, err := svc.Create(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNilModel)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ZeroID", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		input := &models.UserContactRelation{UserID: 0, ContactID: 0}

		result, err := svc.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RelationExists_ReturnsExisting", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}
		existing := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return([]*models.UserContactRelation{existing}, nil)

		result, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, existing, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RelationExists_ButNotFound", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrRelationExists)
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return([]*models.UserContactRelation{
				{UserID: 1, ContactID: 999},
			}, nil)

		result, err := svc.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidForeignKey", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 999}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrDBInvalidForeignKey)

		result, err := svc.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("OtherRepositoryError", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errors.New("db error"))

		result, err := svc.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "erro ao criar")
		mockRepo.AssertExpectations(t)
	})

	t.Run("RelationExists_GetAllError", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		svc := NewUserContactRelationService(mockRepo)

		input := &models.UserContactRelation{UserID: 1, ContactID: 2}

		// Simula que a criação retorna ErrRelationExists
		mockRepo.On("Create", mock.Anything, input).Return(nil, errMsg.ErrRelationExists)
		// Simula erro ao buscar relações existentes
		mockRepo.On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return(nil, errors.New("db error"))

		result, err := svc.Create(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "erro ao verificar relação") // ou parte da mensagem que você usa no fmt.Errorf
		mockRepo.AssertExpectations(t)
	})

}

func TestUserContactRelationServices_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		userID := int64(1)
		contactID := int64(2)

		mockRepo.On("Delete", mock.Anything, userID, contactID).
			Return(nil)

		err := service.Delete(context.Background(), userID, contactID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - userID inválido", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		err := service.Delete(context.Background(), 0, 2)

		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro - contactID inválido", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		err := service.Delete(context.Background(), 1, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro - relação não encontrada", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		userID := int64(1)
		contactID := int64(2)

		mockRepo.On("Delete", mock.Anything, userID, contactID).
			Return(errMsg.ErrNotFound)

		err := service.Delete(context.Background(), userID, contactID)

		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado do repositório", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		userID := int64(1)
		contactID := int64(2)
		expectedErr := errors.New("db error")

		mockRepo.On("Delete", mock.Anything, userID, contactID).
			Return(expectedErr)

		err := service.Delete(context.Background(), userID, contactID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserContactRelationServices_DeleteAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		userID := int64(1)

		mockRepo.On("DeleteAll", mock.Anything, userID).
			Return(nil)

		err := service.DeleteAll(context.Background(), userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro - userID inválido", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		err := service.DeleteAll(context.Background(), 0)

		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "DeleteAll")
	})

	t.Run("erro inesperado do repositório", func(t *testing.T) {
		mockRepo := new(mockContRel.MockUserContactRelation)
		service := NewUserContactRelationService(mockRepo)

		userID := int64(1)
		expectedErr := errors.New("db error")

		mockRepo.On("DeleteAll", mock.Anything, userID).
			Return(expectedErr)

		err := service.DeleteAll(context.Background(), userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}
