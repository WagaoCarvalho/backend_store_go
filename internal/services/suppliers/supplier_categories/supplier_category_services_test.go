package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_categories"
)

// Mock repository
type MockSupplierCategoryRepo struct {
	mock.Mock
}

func (m *MockSupplierCategoryRepo) Create(ctx context.Context, category *models.SupplierCategory) (int64, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierCategoryRepo) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryRepo) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryRepo) Update(ctx context.Context, category *models.SupplierCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockSupplierCategoryRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateSupplierCategory_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	category := &models.SupplierCategory{Name: "Alimentos"}

	mockRepo.On("Create", mock.Anything, category).Return(int64(1), nil)

	id, err := service.Create(context.Background(), category)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestCreateSupplierCategory_InvalidName(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	category := &models.SupplierCategory{Name: " "}

	id, err := service.Create(context.Background(), category)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetSupplierCategoryByID_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	expected := &models.SupplierCategory{ID: 1, Name: "Eletrônicos"}
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	result, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetSupplierCategoryByID_InvalidID(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	result, err := service.GetByID(context.Background(), -1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestUpdateSupplierCategory_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	category := &models.SupplierCategory{ID: 1, Name: "Atualizada"}
	mockRepo.On("Update", mock.Anything, category).Return(nil)

	err := service.Update(context.Background(), category)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSupplierCategory_MissingID(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	category := &models.SupplierCategory{ID: 0, Name: "Nome"}

	err := service.Update(context.Background(), category)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateSupplierCategory_InvalidName(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	category := &models.SupplierCategory{ID: 1, Name: ""}

	err := service.Update(context.Background(), category)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestDeleteSupplierCategory_Success(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := service.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSupplierCategory_InvalidID(t *testing.T) {
	mockRepo := new(MockSupplierCategoryRepo)
	service := services.NewSupplierCategoryService(mockRepo)

	err := service.Delete(context.Background(), 0)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Delete")
}
