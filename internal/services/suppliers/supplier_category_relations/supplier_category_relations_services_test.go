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

// Mock do repositório
type MockSupplierCategoryRelationRepo struct {
	mock.Mock
}

func (m *MockSupplierCategoryRelationRepo) Create(ctx context.Context, rel *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, rel)
	return args.Get(0).(*models.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*models.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationRepo) Delete(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationRepo) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationRepo) HasSupplierCategoryRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}

func TestCreate(t *testing.T) {
	t.Run("IDs inválidos", func(t *testing.T) {
		s := service.NewSupplierCategoryRelationService(nil)

		_, err := s.Create(context.Background(), 0, 1)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)

		_, err = s.Create(context.Background(), 1, 0)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	})

	t.Run("falha ao verificar se relação existe", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).
			Return(false, errors.New("erro ao verificar"))

		_, err := s.Create(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrCheckRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("relação já existe", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).
			Return(true, nil)

		_, err := s.Create(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrRelationExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao criar a relação", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
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
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).
			Return(false, nil)

		expected := &models.SupplierCategoryRelations{ID: 1, SupplierID: 1, CategoryID: 2}
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
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		expected := []*models.SupplierCategoryRelations{
			{ID: 1, SupplierID: 1, CategoryID: 2},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := s.GetBySupplierId(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
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
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		expected := []*models.SupplierCategoryRelations{
			{ID: 1, SupplierID: 1, CategoryID: 2},
		}

		mockRepo.On("GetByCategoryID", mock.Anything, int64(2)).Return(expected, nil)

		result, err := s.GetByCategoryId(context.Background(), 2)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("GetByCategoryID", mock.Anything, int64(2)).Return(([]*models.SupplierCategoryRelations)(nil), errors.New("erro no banco"))

		result, err := s.GetByCategoryId(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestHasRelation(t *testing.T) {
	t.Run("IDs inválidos", func(t *testing.T) {
		s := service.NewSupplierCategoryRelationService(nil)

		result, err := s.HasRelation(context.Background(), 0, 2)

		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
		assert.False(t, result)
	})

	t.Run("relação existe", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).Return(true, nil)

		result, err := s.HasRelation(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao verificar existência", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("HasSupplierCategoryRelation", mock.Anything, int64(1), int64(2)).Return(false, errors.New("erro no banco"))

		result, err := s.HasRelation(context.Background(), 1, 2)

		assert.Error(t, err)
		assert.False(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteById(t *testing.T) {
	t.Run("IDs inválidos", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		err := s.DeleteById(context.Background(), 0, 2)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)

		err = s.DeleteById(context.Background(), 1, -5)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	})

	t.Run("remoção com sucesso", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

		err := s.DeleteById(context.Background(), 1, 2)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("relação não encontrada", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(repository.ErrRelationNotFound)

		err := s.DeleteById(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrRelationNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro inesperado ao deletar", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("erro de banco"))

		err := s.DeleteById(context.Background(), 1, 2)
		assert.ErrorIs(t, err, service.ErrDeleteRelation)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteAllBySupplierId(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		err := s.DeleteAllBySupplierId(context.Background(), 0)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	})

	t.Run("remoção com sucesso", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("DeleteAllBySupplierId", mock.Anything, int64(1)).Return(nil)

		err := s.DeleteAllBySupplierId(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(MockSupplierCategoryRelationRepo)
		s := service.NewSupplierCategoryRelationService(mockRepo)

		mockRepo.On("DeleteAllBySupplierId", mock.Anything, int64(99)).Return(errors.New("falha ao deletar"))

		err := s.DeleteAllBySupplierId(context.Background(), 99)
		assert.ErrorContains(t, err, "falha ao deletar")
		mockRepo.AssertExpectations(t)
	})
}
