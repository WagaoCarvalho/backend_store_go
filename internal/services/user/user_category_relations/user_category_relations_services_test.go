package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repoMocks "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
)

// Mock Repository
type MockUserCategoryRelationRepo struct {
	mock.Mock
}

func (m *MockUserCategoryRelationRepo) Create(ctx context.Context, rel models.UserCategoryRelations) (models.UserCategoryRelations, error) {
	args := m.Called(ctx, rel)
	return args.Get(0).(models.UserCategoryRelations), args.Error(1)
}

func (m *MockUserCategoryRelationRepo) GetByUserID(ctx context.Context, userID int64) ([]models.UserCategoryRelations, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.UserCategoryRelations), args.Error(1)
}

func (m *MockUserCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]models.UserCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]models.UserCategoryRelations), args.Error(1)
}

func (m *MockUserCategoryRelationRepo) Delete(ctx context.Context, userID, categoryID int64) error {
	args := m.Called(ctx, userID, categoryID)
	return args.Error(0)
}

func (m *MockUserCategoryRelationRepo) DeleteAll(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	input := models.UserCategoryRelations{UserID: 1, CategoryID: 2}
	expected := input

	mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

	result, err := service.Create(context.Background(), 1, 2)

	assert.NoError(t, err)
	assert.Equal(t, expected, *result)
	mockRepo.AssertExpectations(t)
}

func TestCreate_InvalidIDs(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	_, err := service.Create(context.Background(), 0, 1)
	assert.ErrorIs(t, err, ErrInvalidUserID)

	_, err = service.Create(context.Background(), 1, 0)
	assert.ErrorIs(t, err, ErrInvalidCategoryID)
}

func TestCreate_AlreadyExists_ReturnsExisting(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	existing := models.UserCategoryRelations{UserID: 1, CategoryID: 2}
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(models.UserCategoryRelations{}, repoMocks.ErrRelationExists)
	mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]models.UserCategoryRelations{existing}, nil)

	result, err := service.Create(context.Background(), 1, 2)

	assert.NoError(t, err)
	assert.Equal(t, existing, *result)
	mockRepo.AssertExpectations(t)
}

func TestCreate_Error(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(models.UserCategoryRelations{}, errors.New("db error"))

	_, err := service.Create(context.Background(), 1, 2)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAll_Success(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	expected := []models.UserCategoryRelations{{UserID: 1, CategoryID: 2}}
	mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(expected, nil)

	result, err := service.GetAll(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAll_InvalidUserID(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	_, err := service.GetAll(context.Background(), 0)
	assert.ErrorIs(t, err, ErrInvalidUserID)
}

func TestGetRelations_Success(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	expected := []models.UserCategoryRelations{{UserID: 1, CategoryID: 2}}
	mockRepo.On("GetByCategoryID", mock.Anything, int64(2)).Return(expected, nil)

	result, err := service.GetRelations(context.Background(), 2)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetRelations_InvalidCategoryID(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	_, err := service.GetRelations(context.Background(), 0)
	assert.ErrorIs(t, err, ErrInvalidCategoryID)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil)

	err := service.Delete(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_InvalidIDs(t *testing.T) {
	service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

	err := service.Delete(context.Background(), 0, 1)
	assert.ErrorIs(t, err, ErrInvalidUserID)

	err = service.Delete(context.Background(), 1, 0)
	assert.ErrorIs(t, err, ErrInvalidCategoryID)
}

func TestDeleteAll_Success(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	mockRepo.On("DeleteAll", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteAll(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAll_InvalidUserID(t *testing.T) {
	service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

	err := service.DeleteAll(context.Background(), 0)
	assert.ErrorIs(t, err, ErrInvalidUserID)
}

func TestGetByCategoryID_SuccessTrue(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	userID := int64(1)
	categoryID := int64(10)

	relations := []models.UserCategoryRelations{
		{UserID: userID, CategoryID: 5},
		{UserID: userID, CategoryID: categoryID}, // match
	}

	mockRepo.On("GetByUserID", mock.Anything, userID).Return(relations, nil)

	result, err := service.GetByCategoryID(context.Background(), userID, categoryID)

	assert.NoError(t, err)
	assert.True(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByCategoryID_SuccessFalse(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	userID := int64(1)
	categoryID := int64(99)

	relations := []models.UserCategoryRelations{
		{UserID: userID, CategoryID: 1},
		{UserID: userID, CategoryID: 2},
	}

	mockRepo.On("GetByUserID", mock.Anything, userID).Return(relations, nil)

	result, err := service.GetByCategoryID(context.Background(), userID, categoryID)

	assert.NoError(t, err)
	assert.False(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByCategoryID_InvalidUserID(t *testing.T) {
	service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

	result, err := service.GetByCategoryID(context.Background(), 0, 1)

	assert.ErrorIs(t, err, ErrInvalidUserID)
	assert.False(t, result)
}

func TestGetByCategoryID_InvalidCategoryID(t *testing.T) {
	service := NewUserCategoryRelationServices(new(MockUserCategoryRelationRepo))

	result, err := service.GetByCategoryID(context.Background(), 1, 0)

	assert.ErrorIs(t, err, ErrInvalidCategoryID)
	assert.False(t, result)
}

func TestGetByCategoryID_RepoError(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepo)
	service := NewUserCategoryRelationServices(mockRepo)

	mockRepo.On("GetByUserID", mock.Anything, int64(1)).
		Return([]models.UserCategoryRelations(nil), errors.New("erro inesperado"))

	result, err := service.GetByCategoryID(context.Background(), 1, 2)

	assert.Error(t, err)
	assert.False(t, result)
	mockRepo.AssertExpectations(t)
}
