package services

import (
	"context"
	"fmt"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	expectedUsers := []models_user.User{
		{UID: 1, Username: "user1", Email: "user1@example.com", Status: true},
		{UID: 2, Username: "user2", Email: "user2@example.com", Status: false},
	}
	mockRepo.On("GetUsers", mock.Anything).Return(expectedUsers, nil)

	userService := NewUserService(mockRepo)
	users, err := userService.GetUsers(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserById(t *testing.T) {
	mockRepo := new(MockUserRepository)
	expectedUser := models_user.User{UID: 1, Username: "user1", Email: "user1@example.com", Status: true}
	mockRepo.On("GetUserById", mock.Anything, int64(1)).Return(expectedUser, nil)

	userService := NewUserService(mockRepo)
	user, err := userService.GetUserById(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserById_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetUserById", mock.Anything, int64(999)).Return(models_user.User{}, fmt.Errorf("usuário não encontrado"))

	userService := NewUserService(mockRepo)
	user, err := userService.GetUserById(context.Background(), 999)

	assert.ErrorContains(t, err, "usuário não encontrado")
	assert.Equal(t, models_user.User{}, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	expectedUser := models_user.User{UID: 1, Username: "user1", Email: "user1@example.com", Status: true}
	mockRepo.On("GetUserByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

	userService := NewUserService(mockRepo)
	user, err := userService.GetUserByEmail(context.Background(), "user1@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByEmail_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(models_user.User{}, fmt.Errorf("usuário não encontrado"))

	userService := NewUserService(mockRepo)
	user, err := userService.GetUserByEmail(context.Background(), "notfound@example.com")

	assert.ErrorContains(t, err, "usuário não encontrado")
	assert.Equal(t, models_user.User{}, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	newUser := models_user.User{Username: "newuser", Email: "newuser@example.com", Status: true}
	categoryID := int64(1)
	newAddress := models_address.Address{
		Street:     "Rua Teste",
		City:       "Cidade Teste",
		State:      "Estado Teste",
		Country:    "Brasil",
		PostalCode: "12345-678",
	}

	mockRepo.On("CreateUser", mock.Anything, newUser, categoryID, newAddress).Return(newUser, nil)

	userService := NewUserService(mockRepo)
	createdUser, err := userService.CreateUser(context.Background(), newUser, categoryID, newAddress)

	assert.NoError(t, err)
	assert.Equal(t, newUser, createdUser)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	newUser := models_user.User{Username: "failuser", Email: "failuser@example.com", Status: true}
	categoryID := int64(1)
	newAddress := models_address.Address{
		Street:     "Rua Teste",
		City:       "Cidade Teste",
		State:      "Estado Teste",
		Country:    "Brasil",
		PostalCode: "12345-678",
	}

	mockRepo.On("CreateUser", mock.Anything, newUser, categoryID, newAddress).Return(models_user.User{}, fmt.Errorf("erro ao criar usuário"))

	userService := NewUserService(mockRepo)
	createdUser, err := userService.CreateUser(context.Background(), newUser, categoryID, newAddress)

	assert.ErrorContains(t, err, "erro ao criar usuário")
	assert.Equal(t, models_user.User{}, createdUser)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	updatedUser := models_user.User{UID: 1, Username: "updateduser", Email: "updated@example.com", Status: true}
	mockRepo.On("UpdateUser", mock.Anything, updatedUser).Return(updatedUser, nil)

	userService := NewUserService(mockRepo)
	resultUser, err := userService.UpdateUser(context.Background(), updatedUser)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, resultUser)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	updatedUser := models_user.User{UID: 1, Username: "failuser", Email: "failuser@example.com", Status: true}
	mockRepo.On("UpdateUser", mock.Anything, updatedUser).Return(models_user.User{}, fmt.Errorf("erro ao atualizar usuário"))

	userService := NewUserService(mockRepo)
	resultUser, err := userService.UpdateUser(context.Background(), updatedUser)

	assert.ErrorContains(t, err, "erro ao atualizar usuário")
	assert.Equal(t, models_user.User{}, resultUser)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUserById(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("DeleteUserById", mock.Anything, int64(1)).Return(nil)

	userService := NewUserService(mockRepo)
	err := userService.DeleteUserById(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUserById_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("DeleteUserById", mock.Anything, int64(999)).Return(fmt.Errorf("usuário não encontrado"))

	userService := NewUserService(mockRepo)
	err := userService.DeleteUserById(context.Background(), 999)

	assert.ErrorContains(t, err, "usuário não encontrado")
	mockRepo.AssertExpectations(t)
}
