package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	repo "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func Test_SupplierContactRelationServices_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("Erro: relação nula", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		result, err := service.Create(ctx, nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrNilModel)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		relation := &models.SupplierContactRelation{SupplierID: 0, ContactID: 2}
		result, err := service.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		relation := &models.SupplierContactRelation{SupplierID: 1, ContactID: 0}
		result, err := service.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Sucesso: relação criada com sucesso", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		expected := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}
		mockRepo.On("Create", ctx, expected).Return(expected, nil).Once()

		result, err := service.Create(ctx, expected)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: relação já existe e GetAllRelationsBySupplierID retorna erro", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		relation := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, relation).Return(nil, err_msg.ErrRelationExists).Once()
		mockRepo.On("GetAllRelationsBySupplierID", ctx, int64(1)).
			Return(nil, errors.New("db error")).Once()

		result, err := service.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrRelationCheck)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: relação já existe e encontrada na lista", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		relation := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, relation).Return(nil, err_msg.ErrRelationExists).Once()
		mockRepo.On("GetAllRelationsBySupplierID", ctx, int64(1)).
			Return([]*models.SupplierContactRelation{{SupplierID: 1, ContactID: 2}}, nil).Once()

		result, err := service.Create(ctx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.ContactID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: relação já existe mas não encontrada", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		relation := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, relation).Return(nil, err_msg.ErrRelationExists).Once()
		mockRepo.On("GetAllRelationsBySupplierID", ctx, int64(1)).
			Return([]*models.SupplierContactRelation{{SupplierID: 1, ContactID: 3}}, nil).Once()

		result, err := service.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: chave estrangeira inválida", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		relation := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}

		mockRepo.On("Create", ctx, relation).Return(nil, err_msg.ErrDBInvalidForeignKey).Once()

		result, err := service.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrDBInvalidForeignKey)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: erro genérico no Create", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		relation := &models.SupplierContactRelation{SupplierID: 1, ContactID: 2}
		dbErr := errors.New("db error")

		mockRepo.On("Create", ctx, relation).Return(nil, dbErr).Once()

		result, err := service.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, err_msg.ErrCreate)
		mockRepo.AssertExpectations(t)
	})
}

func Test_SupplierContactRelationServices_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := service.Delete(context.Background(), 1, 2)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("db error"))

		err := service.Delete(context.Background(), 1, 2)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		err := service.Delete(context.Background(), 0, 2)

		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: contactID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		err := service.Delete(context.Background(), 1, 0)

		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})
}

func Test_SupplierContactRelationServices_DeleteAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil)

		err := service.DeleteAll(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(errors.New("db error"))

		err := service.DeleteAll(context.Background(), 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro: supplierID inválido", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierContactRelation)
		service := NewSupplierContactRelation(mockRepo)

		err := service.DeleteAll(context.Background(), 0)

		assert.ErrorIs(t, err, err_msg.ErrZeroID)
		mockRepo.AssertExpectations(t)
	})
}
