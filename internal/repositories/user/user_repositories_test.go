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

// Testes para GetUsers
func TestGetUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	expectedUsers := []models.User{
		{UID: 1, Username: "user1", Email: "user1@example.com", Password: "hash1", Status: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{UID: 2, Username: "user2", Email: "user2@example.com", Password: "hash2", Status: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockRepo.On("GetUsers", mock.Anything).Return(expectedUsers, nil)

	users, err := mockRepo.GetUsers(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	mockRepo.AssertCalled(t, "GetUsers", mock.Anything)
}

// Testes para GetUserById
func TestGetUserById_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	expectedUser := models.User{UID: 1, Username: "user1", Email: "user1@example.com", Password: "hash1", Status: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mockRepo.On("GetUserById", mock.Anything, int64(1)).Return(expectedUser, nil)

	user, err := mockRepo.GetUserById(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(1))
}

func TestGetUserById_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserById", mock.Anything, int64(999)).Return(models.User{}, fmt.Errorf("usuário não encontrado"))

	user, err := mockRepo.GetUserById(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, "usuário não encontrado", err.Error())
	assert.Equal(t, models.User{}, user)

	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(999))
}

// Testes para CreateUser
func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	inputUser := models.User{Username: "user1", Email: "user1@example.com", Password: "plaintextpassword", Status: true}
	expectedUser := inputUser
	expectedUser.UID = 1
	expectedUser.Password = "hashedpassword"

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(expectedUser, nil)

	user, err := mockRepo.CreateUser(context.Background(), inputUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(models.User{}, fmt.Errorf("email já cadastrado"))

	_, err := mockRepo.CreateUser(context.Background(), models.User{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email já cadastrado")

	mockRepo.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

// Testes para UpdateUser
func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	inputUser := models.User{UID: 1, Username: "updatedUser", Email: "updated@example.com", Password: "newpassword", Status: false}
	updatedUser := inputUser
	updatedUser.UpdatedAt = time.Now()

	mockRepo.On("UpdateUser", mock.Anything, inputUser).Return(updatedUser, nil)

	user, err := mockRepo.UpdateUser(context.Background(), inputUser)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, user)

	mockRepo.AssertCalled(t, "UpdateUser", mock.Anything, inputUser)
}

func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(models.User{}, fmt.Errorf("usuário não encontrado"))

	_, err := mockRepo.UpdateUser(context.Background(), models.User{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usuário não encontrado")

	mockRepo.AssertCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

// Testes para DeleteUserById
func TestDeleteUserById_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	userID := int64(1)

	mockRepo.On("DeleteUserById", mock.Anything, userID).Return(nil)

	err := mockRepo.DeleteUserById(context.Background(), userID)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "DeleteUserById", mock.Anything, userID)
}

func TestDeleteUserById_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("DeleteUserById", mock.Anything, int64(999)).Return(fmt.Errorf("usuário não encontrado"))

	err := mockRepo.DeleteUserById(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, "usuário não encontrado", err.Error())

	mockRepo.AssertCalled(t, "DeleteUserById", mock.Anything, int64(999))
}

func TestDeleteUserById_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("DeleteUserById", mock.Anything, int64(1)).Return(fmt.Errorf("erro ao deletar usuário"))

	err := mockRepo.DeleteUserById(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao deletar usuário")

	mockRepo.AssertCalled(t, "DeleteUserById", mock.Anything, int64(1))
}
