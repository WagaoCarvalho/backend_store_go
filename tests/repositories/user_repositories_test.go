package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUsers(t *testing.T) {
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

	users, err := mockRepo.GetUsers(context.Background())

	assert.NoError(t, err)

	assert.Equal(t, expectedUsers, users)

	mockRepo.AssertCalled(t, "GetUsers", mock.Anything)
}

func TestGetUserById(t *testing.T) {
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

	user, err := mockRepo.GetUserById(context.Background(), 1)

	assert.NoError(t, err)

	assert.Equal(t, expectedUser, user)

	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(1))
}

func TestGetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)

	expectedUser := models.User{
		UID:       1,
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "hashedpassword",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

	user, err := mockRepo.GetUserByEmail(context.Background(), "user1@example.com")

	assert.NoError(t, err)

	assert.Equal(t, expectedUser, user)

	mockRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, "user1@example.com")
}

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").Return(models.User{}, fmt.Errorf("usuário não encontrado"))

	user, err := mockRepo.GetUserByEmail(context.Background(), "nonexistent@example.com")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usuário não encontrado")
	assert.Equal(t, models.User{}, user)

	mockRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, "nonexistent@example.com")
}
