package services

import (
	"context"
	"errors"
	"testing"
	"time"

	repo "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SupplierContactRelationServices_GetAllRelationsBySupplierID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		expected := []*models.SupplierContactRelation{
			{SupplierID: 1, ContactID: 2, CreatedAt: time.Now()},
		}

		mockRepo.On("GetAllRelationsBySupplierID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := service.GetAllRelationsBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		result, err := service.GetAllRelationsBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: falha no repositório", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("GetAllRelationsBySupplierID", mock.Anything, int64(1)).Return(nil, errors.New("db error"))

		result, err := service.GetAllRelationsBySupplierID(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrGet)
		mockRepo.AssertExpectations(t)
	})
}

func Test_SupplierContactRelationServices_HasSupplierContactRelation(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("HasSupplierContactRelation", mock.Anything, int64(1), int64(2)).Return(true, nil)

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotExists", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("HasSupplierContactRelation", mock.Anything, int64(1), int64(3)).Return(false, nil)

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 3)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		exists, err := service.HasSupplierContactRelation(context.Background(), 0, 2)

		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 0)

		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: falha no repositório", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("HasSupplierContactRelation", mock.Anything, int64(1), int64(2)).Return(false, errors.New("db error"))

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 2)

		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrGet)
		mockRepo.AssertExpectations(t)
	})
}
