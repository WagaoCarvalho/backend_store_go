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
	repositories_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Create(t *testing.T) {

	setup := func() (*MockUserRepository, *MockAddressRepository, *MockContactRepository, *MockUserCategoryRelationRepositories, UserService) {
		mockUserRepo := new(MockUserRepository)
		mockAddressRepo := new(MockAddressRepository)
		mockContactRepo := new(MockContactRepository)
		mockRelationRepo := new(MockUserCategoryRelationRepositories)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, userService
	}

	t.Run("sucesso ao criar usuário com todos os dados", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, userService := setup()

		newUser := models_user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Status:   true,
		}
		categoryIDs := []int64{1, 2}
		newAddress := models_address.Address{
			Street:     "Rua Teste",
			City:       "São Paulo",
			PostalCode: "12345-678",
		}
		newContact := models_contact.Contact{
			ContactName: "Contato Teste",
			Email:       "contato@test.com",
			Phone:       "123456789",
		}

		createdUser := models_user.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
			Status:   true,
		}
		mockUserRepo.On("Create", mock.Anything, &newUser).Return(createdUser, nil)
		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(models_address.Address{ID: ptrInt64(1)}, nil)
		mockContactRepo.On("Create", mock.Anything, mock.Anything).Return(&models_contact.Contact{ID: ptrInt64(1)}, nil)
		mockRelationRepo.On("Create", mock.Anything, mock.Anything).Return(&models_user_categories.UserCategory{}, nil).Times(len(categoryIDs))

		result, err := userService.Create(context.Background(), &newUser, categoryIDs, &newAddress, &newContact)

		assert.NoError(t, err)
		assert.Equal(t, createdUser, result)
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
	})

	t.Run("erro ao criar usuário", func(t *testing.T) {
		mockUserRepo, _, _, _, userService := setup()

		newUser := models_user.User{Email: "test@example.com"}
		mockUserRepo.On("Create", mock.Anything, &newUser).Return(models_user.User{}, errors.New("erro no banco de dados"))

		_, err := userService.Create(context.Background(), &newUser, nil, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar usuário")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("erro ao criar endereço", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, _, _, userService := setup()

		newUser := models_user.User{Email: "test@example.com"}
		createdUser := models_user.User{UID: 1}
		mockUserRepo.On("Create", mock.Anything, &newUser).Return(createdUser, nil)
		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(models_address.Address{}, errors.New("erro no endereço"))

		_, err := userService.Create(context.Background(), &newUser, nil, &models_address.Address{}, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço")
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
	})

	t.Run("erro ao criar contato", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, _, userService := setup()

		newUser := models_user.User{Email: "test@example.com"}
		createdUser := models_user.User{UID: 1}
		mockUserRepo.On("Create", mock.Anything, &newUser).Return(createdUser, nil)
		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(models_address.Address{}, nil)
		mockContactRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("erro no contato"))

		_, err := userService.Create(context.Background(), &newUser, nil, &models_address.Address{}, &models_contact.Contact{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
	})

	t.Run("erro ao criar relação com categoria", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, userService := setup()

		newUser := models_user.User{Email: "test@example.com"}
		createdUser := models_user.User{UID: 1}
		categoryIDs := []int64{1, 2}

		mockUserRepo.On("Create", mock.Anything, &newUser).Return(createdUser, nil)
		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(models_address.Address{}, nil)
		mockContactRepo.On("Create", mock.Anything, mock.Anything).Return(&models_contact.Contact{}, nil)
		mockRelationRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("erro na relação"))

		_, err := userService.Create(context.Background(), &newUser, categoryIDs, &models_address.Address{}, &models_contact.Contact{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar relação com categoria")
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
	})

	t.Run("email inválido", func(t *testing.T) {
		_, _, _, _, userService := setup()

		_, err := userService.Create(context.Background(), &models_user.User{Email: "email-invalido"}, nil, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
	})
}

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

	setupMocks := func() (*MockUserRepository, *MockAddressRepository, *MockContactRepository, *MockUserCategoryRelationRepositories) {
		return new(MockUserRepository),
			new(MockAddressRepository),
			new(MockContactRepository),
			new(MockUserCategoryRelationRepositories)
	}

	t.Run("Deve retornar usuário quando encontrado", func(t *testing.T) {
		mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo := setupMocks()

		expectedUser := models_user.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}

		mockUserRepo.On("GetById", mock.Anything, int64(1)).Return(expectedUser, nil)

		userService := NewUserService(mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo)
		user, err := userService.GetById(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando usuário não existe", func(t *testing.T) {
		mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo := setupMocks()

		mockUserRepo.On("GetById", mock.Anything, int64(999)).Return(
			models_user.User{},
			fmt.Errorf("usuário não encontrado"),
		)

		userService := NewUserService(mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo)
		user, err := userService.GetById(context.Background(), 999)

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Equal(t, models_user.User{}, user)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {

	setup := func() (*MockUserRepository, UserService) {
		mockUserRepo := new(MockUserRepository)

		service := NewUserService(
			mockUserRepo,
			new(MockAddressRepository),
			new(MockContactRepository),
			new(MockUserCategoryRelationRepositories),
		)
		return mockUserRepo, service
	}

	t.Run("Deve retornar usuário quando email existe", func(t *testing.T) {
		mockUserRepo, userService := setup()

		expectedUser := models_user.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}

		mockUserRepo.On("GetByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

		user, err := userService.GetByEmail(context.Background(), "user1@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando email não existe", func(t *testing.T) {
		mockUserRepo, userService := setup()

		mockUserRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(
			models_user.User{},
			fmt.Errorf("usuário não encontrado"),
		)

		user, err := userService.GetByEmail(context.Background(), "notfound@example.com")

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Equal(t, models_user.User{}, user)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {

	setup := func() (*MockUserRepository, UserService) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(
			mockRepo,
			new(MockAddressRepository),
			new(MockContactRepository),
			new(MockUserCategoryRelationRepositories),
		)
		return mockRepo, service
	}

	t.Run("deve atualizar usuário com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		inputUser := &models_user.User{
			UID:      1,
			Username: "user1",
			Email:    "valid@example.com",
			Status:   true,
			Version:  1,
		}

		expectedUser := *inputUser
		expectedUser.Username = "user1-updated"

		mockRepo.On("Update", mock.Anything, *inputUser).Return(expectedUser, nil)

		result, err := service.Update(context.Background(), inputUser)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro para email inválido", func(t *testing.T) {
		_, service := setup()

		invalidUser := &models_user.User{
			Email: "invalid-email",
		}

		result, err := service.Update(context.Background(), invalidUser)

		assert.Error(t, err)
		assert.Equal(t, models_user.User{}, result)
		assert.Contains(t, err.Error(), "email inválido")
	})

	t.Run("deve lidar com conflito de versão", func(t *testing.T) {
		mockRepo, service := setup()

		user := &models_user.User{
			UID:     1,
			Email:   "valid@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, *user).
			Return(models_user.User{}, repositories_user.ErrVersionConflict)

		result, err := service.Update(context.Background(), user)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repositories_user.ErrVersionConflict))
		assert.Equal(t, models_user.User{}, result)
	})

	t.Run("deve lidar com usuário não encontrado", func(t *testing.T) {
		mockRepo, service := setup()

		user := &models_user.User{
			UID:     999,
			Email:   "valid@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, *user).
			Return(models_user.User{}, repositories_user.ErrRecordNotFound)

		result, err := service.Update(context.Background(), user)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repositories_user.ErrRecordNotFound))
		assert.Equal(t, models_user.User{}, result)
	})

	t.Run("deve lidar com outros erros do repositório", func(t *testing.T) {
		mockRepo, service := setup()

		user := &models_user.User{
			UID:     1,
			Email:   "valid@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, *user).
			Return(models_user.User{}, fmt.Errorf("erro no banco de dados"))

		result, err := service.Update(context.Background(), user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar usuário")
		assert.Equal(t, models_user.User{}, result)
	})
}

func TestUserService_Delete(t *testing.T) {

	setup := func() (*MockUserRepository, UserService) {
		mockUserRepo := new(MockUserRepository)
		userService := NewUserService(
			mockUserRepo,
			new(MockAddressRepository),
			new(MockContactRepository),
			new(MockUserCategoryRelationRepositories),
		)
		return mockUserRepo, userService
	}

	t.Run("deve deletar usuário com sucesso", func(t *testing.T) {
		mockUserRepo, service := setup()

		mockUserRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro quando usuário não existe", func(t *testing.T) {
		mockUserRepo, service := setup()

		mockUserRepo.On("Delete", mock.Anything, int64(999)).
			Return(repositories_user.ErrRecordNotFound)

		err := service.Delete(context.Background(), 999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao deletar usuário")
		assert.Contains(t, err.Error(), "registro não encontrado")
		assert.True(t, errors.Is(err, repositories_user.ErrRecordNotFound), "deve envolver o erro original")
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("deve retornar erro genérico do repositório", func(t *testing.T) {
		mockUserRepo, service := setup()

		expectedErr := fmt.Errorf("erro no banco de dados")
		mockUserRepo.On("Delete", mock.Anything, int64(1)).
			Return(expectedErr)

		err := service.Delete(context.Background(), 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao deletar usuário")
		assert.True(t, errors.Is(err, expectedErr), "deve envolver o erro original")
		mockUserRepo.AssertExpectations(t)
	})
}

func ptrInt64(i int64) *int64 {
	return &i
}
