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

func (m *MockSupplierCategoryRelationRepo) DeleteAllBySupplier(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationRepo) CheckIfExists(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	mockRepo.On("CheckIfExists", mock.Anything, int64(1), int64(2)).Return(false, nil)
	expected := &models.SupplierCategoryRelations{ID: 1, SupplierID: 1, CategoryID: 2}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.SupplierCategoryRelations")).Return(expected, nil)

	result, err := s.Create(context.Background(), 1, 2)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreate_InvalidIDs(t *testing.T) {
	s := service.NewSupplierCategoryRelationService(nil)

	_, err := s.Create(context.Background(), 0, 1)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)

	_, err = s.Create(context.Background(), 1, 0)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)
}

func TestCreate_AlreadyExists(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	mockRepo.On("CheckIfExists", mock.Anything, int64(1), int64(2)).Return(true, nil)

	_, err := s.Create(context.Background(), 1, 2)
	assert.ErrorIs(t, err, service.ErrRelationExists)
}

func TestGetBySupplier_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	expected := []*models.SupplierCategoryRelations{{ID: 1, SupplierID: 1, CategoryID: 2}}
	mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return(expected, nil)

	result, err := s.GetBySupplier(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetBySupplier_InvalidID(t *testing.T) {
	s := service.NewSupplierCategoryRelationService(nil)
	result, err := s.GetBySupplier(context.Background(), 0)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	assert.Nil(t, result)
}

func TestGetByCategory_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	expected := []*models.SupplierCategoryRelations{{ID: 1, SupplierID: 1, CategoryID: 2}}
	mockRepo.On("GetByCategoryID", mock.Anything, int64(2)).Return(expected, nil)

	result, err := s.GetByCategory(context.Background(), 2)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetByCategory_InvalidID(t *testing.T) {
	s := service.NewSupplierCategoryRelationService(nil)
	result, err := s.GetByCategory(context.Background(), -1)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	assert.Nil(t, result)
}

func TestHasRelation_Exists(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	mockRepo.On("CheckIfExists", mock.Anything, int64(1), int64(2)).Return(true, nil)

	result, err := s.HasRelation(context.Background(), 1, 2)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestHasRelation_InvalidIDs(t *testing.T) {
	s := service.NewSupplierCategoryRelationService(nil)

	result, err := s.HasRelation(context.Background(), 0, 2)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	assert.False(t, result)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	service := service.NewSupplierCategoryRelationService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

	err := service.Delete(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_InvalidIDs(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	err := s.Delete(context.Background(), 0, 2)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)

	err = s.Delete(context.Background(), 1, -5)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)
}

func TestDelete_RelationNotFound(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(repository.ErrRelationNotFound)

	err := s.Delete(context.Background(), 1, 2)
	assert.ErrorIs(t, err, service.ErrRelationNotFound)
}

func TestDeleteAll_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	mockRepo.On("DeleteAllBySupplier", mock.Anything, int64(1)).Return(nil)

	err := s.DeleteAll(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAll_InvalidSupplierID(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	s := service.NewSupplierCategoryRelationService(mockRepo)

	err := s.DeleteAll(context.Background(), 0)
	assert.ErrorIs(t, err, service.ErrInvalidRelationData)
}

func TestDeleteAll_RepoError(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRelationRepo)
	service := service.NewSupplierCategoryRelationService(mockRepo)

	mockRepo.On("DeleteAllBySupplier", mock.Anything, int64(99)).Return(errors.New("falha ao deletar"))

	err := service.DeleteAll(context.Background(), 99)
	assert.ErrorContains(t, err, "falha ao deletar")
}
