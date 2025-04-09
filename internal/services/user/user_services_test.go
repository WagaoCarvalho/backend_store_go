package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
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
	expectedUser := models_user.User{
		UID:      1,
		Username: "user1",
		Email:    "user1@example.com",
		Status:   true,
		Contact: &models_contact.Contact{
			ID:          1,
			UserID:      ptrInt64(1),
			ContactName: "Contato 1",
			Email:       "contato@example.com",
			Phone:       "123456789",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

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
	newContact := models_contact.Contact{
		ContactName:     "Contato Teste",
		ContactPosition: "Gerente",
		Email:           "contato@example.com",
		Phone:           "1111-2222",
		Cell:            "99999-8888",
		ContactType:     "Pessoal",
	}

	mockRepo.On("CreateUser", mock.Anything, newUser, categoryID, newAddress, newContact).
		Return(newUser, nil)

	userService := NewUserService(mockRepo)
	createdUser, err := userService.CreateUser(context.Background(), newUser, categoryID, newAddress, newContact)

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
	newContact := models_contact.Contact{
		ContactName:     "Contato Erro",
		ContactPosition: "Diretor",
		Email:           "erro@example.com",
		Phone:           "1234-5678",
		Cell:            "99999-7777",
		ContactType:     "Profissional",
	}

	mockRepo.On("CreateUser", mock.Anything, newUser, categoryID, newAddress, newContact).
		Return(models_user.User{}, fmt.Errorf("erro ao criar usuário"))

	userService := NewUserService(mockRepo)
	createdUser, err := userService.CreateUser(context.Background(), newUser, categoryID, newAddress, newContact)

	assert.ErrorContains(t, err, "erro ao criar usuário")
	assert.Equal(t, models_user.User{}, createdUser)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	contact := &models_contact.Contact{
		ID:          1,
		UserID:      ptrInt64(1),
		ContactName: "Contato Atualizado",
		Email:       "contato@exemplo.com",
		Phone:       "999999999",
	}

	updatedUser := models_user.User{
		UID:      1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Status:   true,
		Address: &models_address.Address{
			Street:     "Rua Atualizada",
			City:       "Cidade Nova",
			State:      "Estado Atualizado",
			Country:    "Brasil",
			PostalCode: "00000-000",
		},
		Categories: []models_user_categories.UserCategory{
			{ID: 2, Name: "Admin", Description: "Administração"},
		},
	}

	mockRepo.On("UpdateUser", mock.Anything, updatedUser, contact).Return(updatedUser, nil)

	userService := NewUserService(mockRepo)
	resultUser, err := userService.UpdateUser(context.Background(), updatedUser, contact)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, resultUser)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	contact := &models_contact.Contact{
		ID:          1,
		UserID:      ptrInt64(1),
		ContactName: "Contato com Erro",
		Email:       "erro@exemplo.com",
		Phone:       "888888888",
	}

	updatedUser := models_user.User{
		UID:      1,
		Username: "failuser",
		Email:    "failuser@example.com",
		Status:   true,
	}

	mockRepo.On("UpdateUser", mock.Anything, updatedUser, contact).
		Return(models_user.User{}, fmt.Errorf("erro ao atualizar usuário"))

	userService := NewUserService(mockRepo)
	resultUser, err := userService.UpdateUser(context.Background(), updatedUser, contact)

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

func ptrInt64(i int64) *int64 {
	return &i
}
