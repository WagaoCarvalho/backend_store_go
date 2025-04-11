package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
)

// Mock Repository
type MockSupplierRepo struct {
	mock.Mock
}

func (m *MockSupplierRepo) Create(ctx context.Context, supplier *models.Supplier) (int64, error) {
	args := m.Called(ctx, supplier)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierRepo) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Supplier), args.Error(1)
}

func (m *MockSupplierRepo) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Supplier), args.Error(1)
}

func (m *MockSupplierRepo) Update(ctx context.Context, supplier *models.Supplier) error {
	args := m.Called(ctx, supplier)
	return args.Error(0)
}

func (m *MockSupplierRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateSupplier_Success(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	input := &models.Supplier{Name: "Fornecedor X"}
	mockRepo.On("Create", mock.Anything, input).Return(int64(1), nil)

	resultID, err := service.Create(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), resultID)
	mockRepo.AssertExpectations(t)
}

func TestCreateSupplier_Error(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	input := &models.Supplier{Name: "Fornecedor Y"}
	mockRepo.On("Create", mock.Anything, input).Return(int64(0), errors.New("erro ao criar"))

	id, err := service.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	mockRepo.AssertExpectations(t)
}

func TestCreateSupplier_NameRequired(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	supplier := &models.Supplier{Name: ""}

	id, err := service.Create(context.Background(), supplier)

	assert.Equal(t, int64(0), id)
	assert.EqualError(t, err, "nome do fornecedor é obrigatório")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetSupplierByID_Success(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	expected := &models.Supplier{ID: 1, Name: "Fornecedor Z"}
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	result, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetSupplierByID_NotFound(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return((*models.Supplier)(nil), errors.New("não encontrado"))

	result, err := service.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllSuppliers_Success(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	expected := []*models.Supplier{
		{ID: 1, Name: "Fornecedor 1"},
		{ID: 2, Name: "Fornecedor 2"},
	}

	mockRepo.On("GetAll", mock.Anything).Return(expected, nil)

	result, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllSuppliers_Error(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	mockRepo.On("GetAll", mock.Anything).Return([]*models.Supplier(nil), errors.New("erro de banco"))

	result, err := service.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSupplier_Success(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	supplier := &models.Supplier{ID: 1, Name: "Atualizado"}
	mockRepo.On("Update", mock.Anything, supplier).Return(nil)

	err := service.Update(context.Background(), supplier)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSupplier_Error(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	supplier := &models.Supplier{ID: 1, Name: "Erro"}
	mockRepo.On("Update", mock.Anything, supplier).Return(errors.New("falha ao atualizar"))

	err := service.Update(context.Background(), supplier)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSupplier_NameRequired(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	supplier := &models.Supplier{ID: 1, Name: ""}

	err := service.Update(context.Background(), supplier)

	assert.EqualError(t, err, "nome do fornecedor é obrigatório")
	// Esse assert abaixo garante que o método NUNCA foi chamado
	mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestDeleteSupplier_Success(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := service.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSupplier_Error(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	service := NewSupplierService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(99)).Return(errors.New("erro ao deletar"))

	err := service.Delete(context.Background(), 99)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
