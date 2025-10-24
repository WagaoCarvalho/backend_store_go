package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	repo "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func Test_SupplierContactRelationServices_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		result, created, err := service.Create(ctx, 0, 2)

		assert.Nil(t, result)
		assert.False(t, created)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		result, created, err := service.Create(ctx, 1, 0)

		assert.Nil(t, result)
		assert.False(t, created)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		expected := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, expected).Return(expected, nil)

		result, created, err := service.Create(ctx, 1, 2)

		assert.NoError(t, err)
		assert.True(t, created)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: relação já existe e GetAllRelationsBySupplierID retorna erro", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		expected := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, expected).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsBySupplierID", ctx, int64(1)).Return(nil, errors.New("db error"))

		result, created, err := service.Create(ctx, 1, 2)

		assert.Nil(t, result)
		assert.False(t, created)
		assert.ErrorIs(t, err, err_msg.ErrRelationCheck)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: relação já existe e relação encontrada", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		expected := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, expected).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsBySupplierID", ctx, int64(1)).Return([]*models.SupplierContactRelation{
			{SupplierID: 1, ContactID: 2},
		}, nil)

		result, created, err := service.Create(ctx, 1, 2)

		assert.NoError(t, err)
		assert.False(t, created)
		assert.Equal(t, int64(2), result.ContactID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: relação já existe mas não encontrada", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		expected := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, expected).Return(nil, err_msg.ErrRelationExists)
		mockRepo.On("GetAllRelationsBySupplierID", ctx, int64(1)).Return([]*models.SupplierContactRelation{
			{SupplierID: 1, ContactID: 3},
		}, nil)

		result, created, err := service.Create(ctx, 1, 2)

		assert.Nil(t, result)
		assert.False(t, created)
		assert.ErrorIs(t, err, err_msg.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: foreign key inválida", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		expected := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, expected).Return(nil, err_msg.ErrDBInvalidForeignKey)

		result, created, err := service.Create(ctx, 1, 2)

		assert.Nil(t, result)
		assert.False(t, created)
		assert.ErrorIs(t, err, err_msg.ErrDBInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: erro genérico no Create", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		expected := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}
		dbErr := errors.New("db error")

		mockRepo.On("Create", ctx, expected).Return(nil, dbErr)

		result, created, err := service.Create(ctx, 1, 2)

		assert.Nil(t, result)
		assert.False(t, created)
		assert.ErrorIs(t, err, err_msg.ErrCreate)
		mockRepo.AssertExpectations(t)
	})
}

func Test_SupplierContactRelationServices_GetAllRelationsBySupplierID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
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
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		result, err := service.GetAllRelationsBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: falha no repositório", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
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
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("HasSupplierContactRelation", mock.Anything, int64(1), int64(2)).Return(true, nil)

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotExists", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("HasSupplierContactRelation", mock.Anything, int64(1), int64(3)).Return(false, nil)

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 3)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		exists, err := service.HasSupplierContactRelation(context.Background(), 0, 2)

		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 0)

		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: falha no repositório", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("HasSupplierContactRelation", mock.Anything, int64(1), int64(2)).Return(false, errors.New("db error"))

		exists, err := service.HasSupplierContactRelation(context.Background(), 1, 2)

		assert.False(t, exists)
		assert.ErrorIs(t, err, err_msg.ErrGet)
		mockRepo.AssertExpectations(t)
	})
}

func Test_SupplierContactRelationServices_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := service.Delete(context.Background(), 1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("db error"))

		err := service.Delete(context.Background(), 1, 2)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		err := service.Delete(context.Background(), 0, 2)

		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		err := service.Delete(context.Background(), 1, 0)

		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})
}

func Test_SupplierContactRelationServices_DeleteAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil)

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.DeleteAll(context.Background(), 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelationRepo)
		service := NewSupplierContactRelation(mockRepo)

		err := service.DeleteAll(context.Background(), 0)

		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})
}
