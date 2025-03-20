package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user"
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

func (m *MockUserRepository) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models.User), args.Error(1)
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

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
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

func (m *MockUserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	newUser := models.User{
		Username:  "newuser",
		Email:     "newuser@example.com",
		Password:  "hashedpassword",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("CreateUser", mock.Anything, newUser).Return(newUser, nil)

	userService := services.NewUserService(mockRepo)

	createdUser, err := userService.CreateUser(context.Background(), newUser)

	assert.NoError(t, err)

	assert.Equal(t, newUser, createdUser)

	mockRepo.AssertCalled(t, "CreateUser", mock.Anything, newUser)
}

func TestUserService_CreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	newUser := models.User{
		Username: "failuser",
		Email:    "failuser@example.com",
		Password: "hashedpassword",
		Status:   true,
	}

	mockRepo.On("CreateUser", mock.Anything, newUser).Return(models.User{}, fmt.Errorf("erro ao criar usuário"))

	userService := services.NewUserService(mockRepo)

	createdUser, err := userService.CreateUser(context.Background(), newUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar usuário")

	assert.Equal(t, models.User{}, createdUser)

	mockRepo.AssertCalled(t, "CreateUser", mock.Anything, newUser)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	updatedUser := models.User{
		UID:       1,
		Username:  "updateduser",
		Email:     "updateduser@example.com",
		Password:  "newhashedpassword",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateUser", mock.Anything, updatedUser).Return(updatedUser, nil)

	userService := services.NewUserService(mockRepo)

	resultUser, err := userService.UpdateUser(context.Background(), updatedUser)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, resultUser)

	mockRepo.AssertCalled(t, "UpdateUser", mock.Anything, updatedUser)
}

func TestUserService_UpdateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	updatedUser := models.User{
		UID:      1,
		Username: "failuser",
		Email:    "failuser@example.com",
		Password: "newhashedpassword",
		Status:   true,
	}

	mockRepo.On("UpdateUser", mock.Anything, updatedUser).Return(models.User{}, fmt.Errorf("erro ao atualizar usuário"))

	userService := services.NewUserService(mockRepo)

	resultUser, err := userService.UpdateUser(context.Background(), updatedUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao atualizar usuário")
	assert.Equal(t, models.User{}, resultUser)

	mockRepo.AssertCalled(t, "UpdateUser", mock.Anything, updatedUser)
}
