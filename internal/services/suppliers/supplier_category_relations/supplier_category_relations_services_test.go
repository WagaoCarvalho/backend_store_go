package services_test

import (
	"context"
	"errors"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_category_relations"
	service "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_category_relations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {

	t.Run("dados inválidos - supplierID ou categoryID <= 0", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		tests := []struct {
			supplierID int64
			categoryID int64
		}{
			{0, 1},  // supplierID inválido
			{1, 0},  // categoryID inválido
			{-1, 5}, // supplierID negativo
			{4, -3}, // categoryID negativo
			{0, 0},  // ambos inválidos
		}

		for _, tt := range tests {
			result, err := s.Create(context.Background(), tt.supplierID, tt.categoryID)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, service.ErrInvalidRelationData)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao verificar se relação existe", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).
			Return(false, errors.New("erro ao verificar"))

		_, err := s.Create(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrCheckRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("relação já existe", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).
			Return(true, nil)

		_, err := s.Create(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao criar a relação", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).
			Return(false, nil)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.SupplierCategoryRelations")).
			Return((*models.SupplierCategoryRelations)(nil), errors.New("erro ao criar"))

		_, err := s.Create(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrCreateRelation)
		mockRepo.AssertExpectations(t)
	})

	t.Run("criação bem-sucedida", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).
			Return(false, nil)

		expected := &models.SupplierCategoryRelations{SupplierID: 1, CategoryID: 2}
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.SupplierCategoryRelations")).
			Return(expected, nil)

		result, err := s.Create(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBySupplierId(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		s := service.NewSupplierCategoryRelationService(nil)

		result, err := s.GetBySupplierId(context.Background(), 0)

		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
		assert.Nil(t, result)
	})

	t.Run("busca com sucesso", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		expected := []*models.SupplierCategoryRelations{
			{SupplierID: 1, CategoryID: 2},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := s.GetBySupplierId(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return(([]*models.SupplierCategoryRelations)(nil), errors.New("erro de banco"))

		result, err := s.GetBySupplierId(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByCategoryId(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		s := service.NewSupplierCategoryRelationService(nil)

		result, err := s.GetByCategoryId(context.Background(), -1)

		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
		assert.Nil(t, result)
	})

	t.Run("busca com sucesso", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		expected := []*models.SupplierCategoryRelations{
			{SupplierID: 1, CategoryID: 2},
		}

		mockRepo.On("GetByCategoryID", mock.Anything, int64(2)).Return(expected, nil)

		result, err := s.GetByCategoryId(context.Background(), 2)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("GetByCategoryID", mock.Anything, int64(2)).Return(([]*models.SupplierCategoryRelations)(nil), errors.New("erro no banco"))

		result, err := s.GetByCategoryId(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {

	t.Run("CategoryID inválido (<= 0)", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelations{
			SupplierID: 1,
			CategoryID: 0, // inválido
			Version:    1,
		}

		result, err := svc.Update(context.Background(), relation)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, service.ErrInvalidSupplierCategoryRelationID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SupplierID inválido (<= 0)", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelations{
			SupplierID: 0, // inválido
			CategoryID: 2, // válido
			Version:    1,
		}

		result, err := svc.Update(context.Background(), relation)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ambos SupplierID e CategoryID inválidos (<= 0)", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelations{
			SupplierID: -1,
			CategoryID: 0,
			Version:    1,
		}

		result, err := svc.Update(context.Background(), relation)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, service.ErrInvalidSupplierCategoryRelationID) // <- esta validação vem primeiro
		mockRepo.AssertExpectations(t)
	})

	t.Run("versão ausente", func(t *testing.T) {
		s := service.NewSupplierCategoryRelationService(nil)

		invalid := &models.SupplierCategoryRelations{SupplierID: 1, CategoryID: 2, Version: 0}
		_, err := s.Update(context.Background(), invalid)
		assert.ErrorIs(t, err, service.ErrSupplierCategoryRelationVersionRequired)
	})

	t.Run("relação não encontrada", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelations{SupplierID: 1, CategoryID: 2, Version: 1}

		mockRepo.On("Update", mock.Anything, relation).
			Return(nil, repository.ErrRelationNotFound)

		_, err := s.Update(context.Background(), relation)
		assert.ErrorIs(t, err, service.ErrSupplierCategoryRelationUpdate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelations{SupplierID: 1, CategoryID: 2, Version: 1}

		mockRepo.On("Update", mock.Anything, relation).
			Return(nil, errors.New("falha inesperada"))

		_, err := s.Update(context.Background(), relation)
		assert.ErrorIs(t, err, service.ErrSupplierCategoryRelationUpdate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("atualização bem-sucedida", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelations{SupplierID: 1, CategoryID: 2, Version: 1}

		mockRepo.On("Update", mock.Anything, relation).
			Return(relation, nil)

		result, err := s.Update(context.Background(), relation)
		assert.NoError(t, err)
		assert.Equal(t, relation, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("relação não encontrada (erro propagado)", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelations{
			SupplierID: 1,
			CategoryID: 2,
			Version:    1,
		}

		mockRepo.
			On("Update", mock.Anything, relation).
			Return(nil, service.ErrSupplierCategoryRelationNotFound)

		_, err := s.Update(context.Background(), relation)

		assert.ErrorIs(t, err, service.ErrSupplierCategoryRelationNotFound)
		mockRepo.AssertExpectations(t)
	})

}

func TestHasRelation(t *testing.T) {

	t.Run("dados inválidos - supplierID ou categoryID <= 0", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		tests := []struct {
			supplierID int64
			categoryID int64
		}{
			{0, 1},  // supplierID inválido
			{1, 0},  // categoryID inválido
			{-1, 2}, // supplierID negativo
			{3, -5}, // categoryID negativo
			{0, 0},  // ambos inválidos
		}

		for _, tt := range tests {
			ok, err := svc.HasRelation(context.Background(), tt.supplierID, tt.categoryID)
			assert.False(t, ok)
			assert.ErrorIs(t, err, service.ErrInvalidRelationData)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("relação existe", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).Return(true, nil)

		result, err := s.HasRelation(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao verificar existência", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).Return(false, errors.New("erro no banco"))

		result, err := s.HasRelation(context.Background(), 1, 2)

		assert.Error(t, err)
		assert.False(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteById(t *testing.T) {

	t.Run("dados inválidos - supplierID ou categoryID <= 0", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		tests := []struct {
			supplierID int64
			categoryID int64
		}{
			{0, 1},  // supplierID inválido
			{1, 0},  // categoryID inválido
			{-1, 2}, // supplierID negativo
			{3, -5}, // categoryID negativo
			{0, 0},  // ambos inválidos
		}

		for _, tt := range tests {
			err := svc.DeleteById(context.Background(), tt.supplierID, tt.categoryID)
			assert.ErrorIs(t, err, service.ErrInvalidRelationData)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("remoção com sucesso", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := s.DeleteById(context.Background(), 1, 2)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("relação não encontrada", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(repository.ErrRelationNotFound)

		err := s.DeleteById(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrRelationNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado ao deletar", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("erro de banco"))

		err := s.DeleteById(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrDeleteRelation)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteAllBySupplierId(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		err := s.DeleteAllBySupplierId(context.Background(), 0)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	})

	t.Run("remoção com sucesso", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("DeleteAllBySupplierId", mock.Anything, int64(1)).Return(nil)

		err := s.DeleteAllBySupplierId(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("DeleteAllBySupplierId", mock.Anything, int64(99)).Return(errors.New("falha ao deletar"))

		err := s.DeleteAllBySupplierId(context.Background(), 99)
		assert.ErrorContains(t, err, "falha ao deletar")
		mockRepo.AssertExpectations(t)
	})
}
