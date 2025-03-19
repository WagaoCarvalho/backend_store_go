package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/WagaoCarvalho/backend_store_go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func TestUserService_GetUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)

	expectedUsers := []models.User{
		{
			UID:       1,
			Username:  "user1",
			Email:     "user1@example.com",
			Password:  "hash1",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UID:       2,
			Username:  "user2",
			Email:     "user2@example.com",
			Password:  "hash2",
			Status:    false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("GetUsers", mock.Anything).Return(expectedUsers, nil)

	userService := services.NewUserService(mockRepo)

	users, err := userService.GetUsers(context.Background())

	assert.NoError(t, err)

	assert.Equal(t, expectedUsers, users)

	mockRepo.AssertCalled(t, "GetUsers", mock.Anything)
}

func TestUserService_GetUserById(t *testing.T) {
	mockRepo := new(MockUserRepository)

	expectedUser := models.User{
		UID:       1,
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "hash1",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetUserById", mock.Anything, int64(1)).Return(expectedUser, nil)

	userService := services.NewUserService(mockRepo)

	user, err := userService.GetUserById(context.Background(), 1)

	assert.NoError(t, err)

	assert.Equal(t, expectedUser, user)

	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(1))
}

func TestUserService_GetUserById_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserById", mock.Anything, int64(999)).Return(models.User{}, fmt.Errorf("usuário não encontrado"))

	userService := services.NewUserService(mockRepo)

	user, err := userService.GetUserById(context.Background(), 999)

	assert.Error(t, err)

	assert.Equal(t, models.User{}, user)

	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(999))
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)

	expectedUser := models.User{
		UID:       1,
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "hash1",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

	userService := services.NewUserService(mockRepo)

	user, err := userService.GetUserByEmail(context.Background(), "user1@example.com")

	assert.NoError(t, err)

	assert.Equal(t, expectedUser, user)

	mockRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, "user1@example.com")
}

func TestUserService_GetUserByEmail_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").Return(models.User{}, fmt.Errorf("usuário não encontrado"))

	userService := services.NewUserService(mockRepo)

	user, err := userService.GetUserByEmail(context.Background(), "nonexistent@example.com")

	assert.Error(t, err)

	assert.Equal(t, models.User{}, user)

	mockRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, "nonexistent@example.com")
}
