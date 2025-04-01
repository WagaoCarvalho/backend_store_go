package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/user_category_relations"
)

// Mock da interface repositories.UserCategoryRelationRepositories
type MockUserCategoryRelationRepository struct {
	mock.Mock
}

func (m *MockUserCategoryRelationRepository) CreateRelation(ctx context.Context, relation models.UserCategoryRelation) (models.UserCategoryRelation, error) {
	args := m.Called(ctx, relation)
	return args.Get(0).(models.UserCategoryRelation), args.Error(1)
}

func (m *MockUserCategoryRelationRepository) GetRelationsByUserID(ctx context.Context, userID int64) ([]models.UserCategoryRelation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.UserCategoryRelation), args.Error(1)
}

func (m *MockUserCategoryRelationRepository) GetRelationsByCategoryID(ctx context.Context, categoryID int64) ([]models.UserCategoryRelation, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]models.UserCategoryRelation), args.Error(1)
}

func (m *MockUserCategoryRelationRepository) DeleteRelation(ctx context.Context, userID, categoryID int64) error {
	args := m.Called(ctx, userID, categoryID)
	return args.Error(0)
}

func (m *MockUserCategoryRelationRepository) DeleteAllUserRelations(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestCreateRelation_Success(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepository)
	service := NewUserCategoryRelationServices(mockRepo)

	expectedRelation := models.UserCategoryRelation{UserID: 1, CategoryID: 2}

	mockRepo.On("CreateRelation", mock.Anything, mock.AnythingOfType("models.UserCategoryRelation")).
		Return(expectedRelation, nil)

	relation, err := service.CreateRelation(context.Background(), 1, 2)

	assert.NoError(t, err)
	assert.NotNil(t, relation)
	assert.Equal(t, int64(1), relation.UserID)
	assert.Equal(t, int64(2), relation.CategoryID)

	mockRepo.AssertExpectations(t)
}

func TestCreateRelation_InvalidUserID(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepository)
	service := NewUserCategoryRelationServices(mockRepo)

	relation, err := service.CreateRelation(context.Background(), -1, 2)

	assert.Error(t, err)
	assert.Nil(t, relation)
	assert.Equal(t, ErrInvalidUserID, err)
}

func TestCreateRelation_InvalidCategoryID(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepository)
	service := NewUserCategoryRelationServices(mockRepo)

	relation, err := service.CreateRelation(context.Background(), 1, -1)

	assert.Error(t, err)
	assert.Nil(t, relation)
	assert.Equal(t, ErrInvalidCategoryID, err)
}

func TestCreateRelation_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepository)
	service := NewUserCategoryRelationServices(mockRepo)

	expectedRelation := models.UserCategoryRelation{UserID: 1, CategoryID: 2}
	existingRelations := []models.UserCategoryRelation{expectedRelation}

	mockRepo.On("CreateRelation", mock.Anything, mock.Anything).
		Return(models.UserCategoryRelation{}, repositories.ErrRelationExists)
	mockRepo.On("GetRelationsByUserID", mock.Anything, int64(1)).
		Return(existingRelations, nil)

	relation, err := service.CreateRelation(context.Background(), 1, 2)

	assert.NoError(t, err)
	assert.NotNil(t, relation)
	assert.Equal(t, int64(1), relation.UserID)
	assert.Equal(t, int64(2), relation.CategoryID)

	mockRepo.AssertExpectations(t)
}

func TestCreateRelation_UnexpectedDBError(t *testing.T) {
	mockRepo := new(MockUserCategoryRelationRepository)
	service := NewUserCategoryRelationServices(mockRepo)

	mockRepo.On("CreateRelation", mock.Anything, mock.Anything).
		Return(models.UserCategoryRelation{}, errors.New("db failure"))

	relation, err := service.CreateRelation(context.Background(), 1, 2)

	assert.Error(t, err)
	assert.Nil(t, relation)
	assert.Contains(t, err.Error(), "erro ao criar relação: db failure")

	mockRepo.AssertExpectations(t)
}
