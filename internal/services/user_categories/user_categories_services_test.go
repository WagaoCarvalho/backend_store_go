package services_test

import (
	"context"
	"errors"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user_categories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserCategoryRepository struct {
	mock.Mock
}

func (m *MockUserCategoryRepository) GetCategories(ctx context.Context) ([]models.UserCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) GetCategoryById(ctx context.Context, id int64) (models.UserCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) CreateCategory(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) UpdateCategory(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) DeleteCategoryById(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserCategoryService_GetCategories(t *testing.T) {
	mockRepo := new(MockUserCategoryRepository)
	expectedCategories := []models.UserCategory{{ID: 1, Name: "Category1", Description: "Desc1"}}
	mockRepo.On("GetCategories", mock.Anything).Return(expectedCategories, nil)

	service := services.NewUserCategoryService(mockRepo)
	categories, err := service.GetCategories(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
	mockRepo.AssertExpectations(t)
}

func TestUserCategoryService_GetCategoryById(t *testing.T) {
	mockRepo := new(MockUserCategoryRepository)
	expectedCategory := models.UserCategory{ID: 1, Name: "Category1", Description: "Desc1"}
	mockRepo.On("GetCategoryById", mock.Anything, int64(1)).Return(expectedCategory, nil)

	service := services.NewUserCategoryService(mockRepo)
	category, err := service.GetCategoryById(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)
	mockRepo.AssertExpectations(t)
}

func TestUserCategoryService_CreateCategory(t *testing.T) {
	mockRepo := new(MockUserCategoryRepository)
	inputCategory := models.UserCategory{Name: "NewCategory", Description: "NewDesc"}
	createdCategory := models.UserCategory{ID: 1, Name: "NewCategory", Description: "NewDesc"}
	mockRepo.On("CreateCategory", mock.Anything, inputCategory).Return(createdCategory, nil)

	service := services.NewUserCategoryService(mockRepo)
	category, err := service.CreateCategory(context.Background(), inputCategory)

	assert.NoError(t, err)
	assert.Equal(t, createdCategory, category)
	mockRepo.AssertExpectations(t)
}

func TestUserCategoryService_UpdateCategory(t *testing.T) {
	mockRepo := new(MockUserCategoryRepository)
	updatedCategory := models.UserCategory{ID: 1, Name: "UpdatedCategory", Description: "UpdatedDesc"}
	mockRepo.On("UpdateCategory", mock.Anything, updatedCategory).Return(updatedCategory, nil)

	service := services.NewUserCategoryService(mockRepo)
	category, err := service.UpdateCategory(context.Background(), updatedCategory)

	assert.NoError(t, err)
	assert.Equal(t, updatedCategory, category)
	mockRepo.AssertExpectations(t)
}

func TestUserCategoryService_DeleteCategoryById(t *testing.T) {
	mockRepo := new(MockUserCategoryRepository)
	mockRepo.On("DeleteCategoryById", mock.Anything, int64(1)).Return(nil)

	service := services.NewUserCategoryService(mockRepo)
	err := service.DeleteCategoryById(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserCategoryService_GetCategoryById_NotFound(t *testing.T) {
	mockRepo := new(MockUserCategoryRepository)
	mockRepo.On("GetCategoryById", mock.Anything, int64(999)).Return(models.UserCategory{}, errors.New("categoria não encontrada"))

	service := services.NewUserCategoryService(mockRepo)
	category, err := service.GetCategoryById(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, "categoria não encontrada", err.Error())
	assert.Equal(t, models.UserCategory{}, category)
	mockRepo.AssertExpectations(t)
}
