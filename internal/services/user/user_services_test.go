package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	user_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetUsers(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	expectedUsers := []models_user.User{
		{UID: 1, Username: "user1", Email: "user1@example.com", Status: true},
		{UID: 2, Username: "user2", Email: "user2@example.com", Status: false},
	}
	mockUserRepo.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)
	users, err := userService.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUserById(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	expectedUser := models_user.User{
		UID:      1,
		Username: "user1",
		Email:    "user1@example.com",
		Status:   true,
	}

	mockUserRepo.On("GetById", mock.Anything, int64(1)).Return(expectedUser, nil)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)
	user, err := userService.GetById(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUserById_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	mockUserRepo.On("GetById", mock.Anything, int64(999)).Return(models_user.User{}, fmt.Errorf("usuário não encontrado"))

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)
	user, err := userService.GetById(context.Background(), 999)

	assert.ErrorContains(t, err, "usuário não encontrado")
	assert.Equal(t, models_user.User{}, user)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	expectedUser := models_user.User{UID: 1, Username: "user1", Email: "user1@example.com", Status: true}
	mockUserRepo.On("GetByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)
	user, err := userService.GetByEmail(context.Background(), "user1@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUserByEmail_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	mockUserRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(models_user.User{}, fmt.Errorf("usuário não encontrado"))

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)
	user, err := userService.GetByEmail(context.Background(), "notfound@example.com")

	assert.ErrorContains(t, err, "usuário não encontrado")
	assert.Equal(t, models_user.User{}, user)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_CreateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	newUser := models_user.User{Username: "newuser", Email: "newuser@example.com", Status: true}
	categoryIDs := []int64{1}

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

	// Configuração dos mocks
	createdUser := models_user.User{UID: 1, Username: "newuser", Email: "newuser@example.com", Status: true}
	mockUserRepo.On("Create", mock.Anything, &newUser).Return(createdUser, nil)
	mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(models_address.Address{}, nil)
	mockContactRepo.On("Create", mock.Anything, mock.Anything).Return(&models_contact.Contact{}, nil)
	mockRelationRepo.On("Create", mock.Anything, mock.Anything).Return(&models_user_categories.UserCategory{}, nil)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)

	resultUser, err := userService.Create(context.Background(), &newUser, categoryIDs, &newAddress, &newContact)

	assert.NoError(t, err)
	assert.Equal(t, createdUser, resultUser)
	mockUserRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
	mockContactRepo.AssertExpectations(t)
	mockRelationRepo.AssertExpectations(t)
}

func TestUserService_Create_ErroAoCriarContato(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	newUser := models_user.User{Username: "newuser", Email: "newuser@example.com", Status: true}
	categoryIDs := []int64{1}
	newAddress := models_address.Address{}
	newContact := models_contact.Contact{}
	createdUser := models_user.User{UID: 1, Username: "newuser", Email: "newuser@example.com", Status: true}

	mockUserRepo.
		On("Create", mock.Anything, &newUser).
		Return(createdUser, nil)

	mockAddressRepo.
		On("Create", mock.Anything, mock.Anything).
		Return(models_address.Address{}, nil)

	// Simula erro ao criar contato
	mockContactRepo.
		On("Create", mock.Anything, mock.Anything).
		Return(nil, errors.New("falha ao criar contato"))

	// Evita panic se for chamado
	mockRelationRepo.
		On("Create", mock.Anything, mock.Anything).
		Maybe().
		Return(&user_category_relations.UserCategoryRelations{}, nil)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)

	_, err := userService.Create(context.Background(), &newUser, categoryIDs, &newAddress, &newContact)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar contato")
	mockUserRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
	mockContactRepo.AssertExpectations(t)
	mockRelationRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything) // Opcional
}

func TestUserService_Create_ErroAoCriarRelacaoCategoria(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	newUser := models_user.User{Username: "newuser", Email: "newuser@example.com", Status: true}
	categoryIDs := []int64{1}
	newAddress := models_address.Address{}
	newContact := models_contact.Contact{}
	createdUser := models_user.User{UID: 1, Username: "newuser", Email: "newuser@example.com", Status: true}

	mockUserRepo.
		On("Create", mock.Anything, &newUser).
		Return(createdUser, nil)

	mockAddressRepo.
		On("Create", mock.Anything, mock.Anything).
		Return(models_address.Address{}, nil)

	mockContactRepo.
		On("Create", mock.Anything, mock.Anything).
		Return(&models_contact.Contact{}, nil)

	mockRelationRepo.
		On("Create", mock.Anything, mock.Anything).
		Return(nil, errors.New("falha ao criar relação"))

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)

	_, err := userService.Create(context.Background(), &newUser, categoryIDs, &newAddress, &newContact)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar relação com categoria ID")
	mockUserRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
	mockContactRepo.AssertExpectations(t)
	mockRelationRepo.AssertExpectations(t)
}

func TestUserService_Create_ErroAoCriarEndereco(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	newUser := models_user.User{Username: "newuser", Email: "newuser@example.com", Status: true}
	categoryIDs := []int64{1}
	newAddress := models_address.Address{}
	newContact := models_contact.Contact{}
	createdUser := models_user.User{UID: 1, Username: "newuser", Email: "newuser@example.com", Status: true}

	mockUserRepo.
		On("Create", mock.Anything, &newUser).
		Return(createdUser, nil)

	mockAddressRepo.
		On("Create", mock.Anything, mock.Anything).
		Return(models_address.Address{}, errors.New("falha ao criar endereço"))

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)

	_, err := userService.Create(context.Background(), &newUser, categoryIDs, &newAddress, &newContact)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar endereço")
	mockUserRepo.AssertExpectations(t)
	mockAddressRepo.AssertExpectations(t)
}

func TestUserService_Create_ErroAoCriarUsuario(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	newUser := models_user.User{Username: "newuser", Email: "newuser@example.com", Status: true}
	categoryIDs := []int64{1}
	newAddress := models_address.Address{}
	newContact := models_contact.Contact{}

	mockUserRepo.
		On("Create", mock.Anything, &newUser).
		Return(models_user.User{}, errors.New("falha ao criar usuário"))

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)

	_, err := userService.Create(context.Background(), &newUser, categoryIDs, &newAddress, &newContact)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar usuário")
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_Create_InvalidEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)

	user := models_user.User{Email: "invalid-email"}
	address := models_address.Address{}
	contact := models_contact.Contact{}
	categoryIDs := []int64{1}

	createdUser, err := userService.Create(context.Background(), &user, categoryIDs, &address, &contact)

	assert.Error(t, err)
	assert.EqualError(t, err, "email inválido")
	assert.Equal(t, models_user.User{}, createdUser)
}

func TestUserService_UpdateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	updatedUser := models_user.User{
		UID:      1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Status:   true,
	}

	mockUserRepo.On("Update", mock.Anything, updatedUser).Return(updatedUser, nil)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)
	resultUser, err := userService.Update(context.Background(), &updatedUser)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, resultUser)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_DeleteUserById(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAddressRepo := new(MockAddressRepository)
	mockContactRepo := new(MockContactRepository)
	mockRelationRepo := new(MockUserCategoryRelationRepositories)

	mockUserRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	userService := NewUserService(mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo)
	err := userService.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func ptrInt64(i int64) *int64 {
	return &i
}
