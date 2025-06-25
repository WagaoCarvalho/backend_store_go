package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_categories_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	addresses_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	repositories_address "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	contact_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	user_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	user_category_relations_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Create(t *testing.T) {

	setup := func() (
		*user_repositories.MockUserRepository,
		*addresses_repositories.MockAddressRepository,
		*contact_repositories.MockContactRepository,
		*user_category_relations_repositories.MockUserCategoryRelationRepo,
		*userService, // <-- trocar aqui de UserService para *userService
	) {
		mockUserRepo := new(user_repositories.MockUserRepository)
		mockAddressRepo := new(addresses_repositories.MockAddressRepository)
		mockContactRepo := new(contact_repositories.MockContactRepository)
		mockRelationRepo := new(user_category_relations_repositories.MockUserCategoryRelationRepo)

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

		newUser := &models_user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Status:   true,
		}
		categoryIDs := []int64{1, 2}
		newAddress := &models_address.Address{
			Street:     "Rua Teste",
			City:       "São Paulo",
			PostalCode: "12345-678",
		}
		newContact := &models_contact.Contact{
			ContactName: "Contato Teste",
			Email:       "contato@test.com",
			Phone:       "123456789",
		}

		createdUser := &models_user.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
			Status:   true,
		}
		mockUserRepo.
			On("Create", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				args.Get(1).(*models_user.User).UID = 1
			}).
			Return(&models_user.User{
				UID:      1,
				Username: "testuser",
				Email:    "test@example.com",
				Status:   true,
			}, nil)

		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(&models_address.Address{ID: 1}, nil)
		mockContactRepo.On("Create", mock.Anything, mock.Anything).Return(&models_contact.Contact{ID: 1}, nil)
		mockRelationRepo.On("Create", mock.Anything, mock.Anything).Return(&models_user_categories_relations.UserCategoryRelations{}, nil).Times(len(categoryIDs))

		result, err := userService.Create(context.Background(), newUser, categoryIDs, newAddress, newContact)

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
		mockUserRepo.
			On("Create", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				args.Get(1).(*models_user.User).UID = 1
			}).
			Return(&models_user.User{
				UID:      1,
				Username: "testuser",
				Email:    "test@example.com",
				Status:   true,
			}, nil)

		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(&models_address.Address{}, errors.New("erro no endereço"))
		_, err := userService.Create(context.Background(), &newUser, nil, &models_address.Address{}, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço")
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
	})

	t.Run("erro ao criar contato", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, _, userService := setup()

		newUser := models_user.User{Email: "test@example.com"}
		mockUserRepo.
			On("Create", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				args.Get(1).(*models_user.User).UID = 1
			}).
			Return(&models_user.User{
				UID:      1,
				Username: "testuser",
				Email:    "test@example.com",
				Status:   true,
			}, nil)
		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(&models_address.Address{}, nil)
		mockContactRepo.On("Create", mock.Anything, mock.Anything).Return((*models_contact.Contact)(nil), errors.New("erro no contato"))

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
		categoryIDs := []int64{1, 2}

		mockUserRepo.
			On("Create", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				args.Get(1).(*models_user.User).UID = 1
			}).
			Return(&models_user.User{
				UID:      1,
				Username: "testuser",
				Email:    "test@example.com",
				Status:   true,
			}, nil)
		mockAddressRepo.On("Create", mock.Anything, mock.Anything).Return(&models_address.Address{}, nil)
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

	t.Run("usuário criado é nulo", func(t *testing.T) {
		mockUserRepo, _, _, _, userService := setup()

		user := &models_user.User{
			Email: "valid@email.com",
		}

		mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil, nil)

		_, err := userService.Create(context.Background(), user, nil, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "usuário criado é nulo")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("email inválido", func(t *testing.T) {
		_, _, _, _, userService := setup()

		_, err := userService.Create(context.Background(), &models_user.User{Email: "email-invalido"}, nil, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
	})
}

func TestUserService_GetUsers(t *testing.T) {
	mockUserRepo := new(user_repositories.MockUserRepository)
	mockAddressRepo := new(addresses_repositories.MockAddressRepository)
	mockContactRepo := new(contact_repositories.MockContactRepository)
	mockRelationRepo := new(user_category_relations_repositories.MockUserCategoryRelationRepo)

	expectedUsers := []*models_user.User{
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

func TestUserService_GetUserByID(t *testing.T) {

	setupMocks := func() (*user_repositories.MockUserRepository, *addresses_repositories.MockAddressRepository, *contact_repositories.MockContactRepository, *user_category_relations_repositories.MockUserCategoryRelationRepo) {
		return new(user_repositories.MockUserRepository),
			new(addresses_repositories.MockAddressRepository),
			new(contact_repositories.MockContactRepository),
			new(user_category_relations_repositories.MockUserCategoryRelationRepo)
	}

	t.Run("Deve retornar usuário quando encontrado", func(t *testing.T) {
		mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo := setupMocks()

		expectedUser := &models_user.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}

		mockUserRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

		userService := NewUserService(mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo)
		user, err := userService.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando usuário não existe", func(t *testing.T) {
		mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo := setupMocks()

		mockUserRepo.On("GetByID", mock.Anything, int64(999)).Return(
			nil, // ponteiro nil
			fmt.Errorf("usuário não encontrado"),
		)

		userService := NewUserService(mockUserRepo, mockAddrRepo, mockContactRepo, mockRelRepo)
		user, err := userService.GetByID(context.Background(), 999)

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Nil(t, user) // agora user deve ser nil
		mockUserRepo.AssertExpectations(t)
	})

}

func TestUserService_GetVersionByID(t *testing.T) {
	t.Run("deve retornar a versão corretamente", func(t *testing.T) {
		mockRepo := new(user_repositories.MockUserRepository)
		service := NewUserService(
			mockRepo,
			new(addresses_repositories.MockAddressRepository),
			new(contact_repositories.MockContactRepository),
			new(user_category_relations_repositories.MockUserCategoryRelationRepo),
		)

		uid := int64(1)
		expectedVersion := int64(5)

		mockRepo.On("GetVersionByID", mock.Anything, uid).Return(expectedVersion, nil).Once()

		version, err := service.GetVersionByID(context.Background(), uid)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo := new(user_repositories.MockUserRepository)
		service := NewUserService(
			mockRepo,
			new(addresses_repositories.MockAddressRepository),
			new(contact_repositories.MockContactRepository),
			new(user_category_relations_repositories.MockUserCategoryRelationRepo),
		)

		uid := int64(999)

		mockRepo.On("GetVersionByID", mock.Anything, uid).
			Return(int64(0), user_repositories.ErrUserNotFound).Once()

		version, err := service.GetVersionByID(context.Background(), uid)

		assert.ErrorIs(t, err, user_repositories.ErrUserNotFound)
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro genérico do repositório", func(t *testing.T) {
		mockRepo := new(user_repositories.MockUserRepository)
		service := NewUserService(
			mockRepo,
			new(addresses_repositories.MockAddressRepository),
			new(contact_repositories.MockContactRepository),
			new(user_category_relations_repositories.MockUserCategoryRelationRepo),
		)

		uid := int64(2)
		repoErr := errors.New("falha no banco")

		mockRepo.On("GetVersionByID", mock.Anything, uid).Return(int64(0), repoErr).Once()

		version, err := service.GetVersionByID(context.Background(), uid)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user: erro ao obter versão")
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {

	setup := func() (*user_repositories.MockUserRepository, UserService) {
		mockUserRepo := new(user_repositories.MockUserRepository)

		service := NewUserService(
			mockUserRepo,
			new(addresses_repositories.MockAddressRepository),
			new(contact_repositories.MockContactRepository),
			new(user_category_relations_repositories.MockUserCategoryRelationRepo),
		)
		return mockUserRepo, service
	}

	t.Run("Deve retornar usuário quando email existe", func(t *testing.T) {
		mockUserRepo, userService := setup()

		expectedUser := &models_user.User{
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
			nil,
			fmt.Errorf("usuário não encontrado"),
		)

		user, err := userService.GetByEmail(context.Background(), "notfound@example.com")

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Nil(t, user)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	setup := func() (*user_repositories.MockUserRepository, *addresses_repositories.MockAddressRepository, UserService) {
		mockUserRepo := new(user_repositories.MockUserRepository)
		mockAddressRepo := new(addresses_repositories.MockAddressRepository)

		UserService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			new(contact_repositories.MockContactRepository),
			new(user_category_relations_repositories.MockUserCategoryRelationRepo),
		)

		return mockUserRepo, mockAddressRepo, UserService
	}

	t.Run("versão inválida", func(t *testing.T) {
		_, _, service := setup()

		user := &models_user.User{
			UID:      1,
			Username: "user1",
			Email:    "valid@example.com",
			Status:   true,
		}

		updated, err := service.Update(context.Background(), user, nil)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, ErrInvalidVersion)
	})

	t.Run("deve atualizar usuário com sucesso", func(t *testing.T) {
		mockRepoUser, _, service := setup()

		inputUser := &models_user.User{
			UID:      1,
			Username: "user1",
			Email:    "valid@example.com",
			Status:   true,
			Version:  1,
		}

		expectedUser := *inputUser
		expectedUser.Username = "user1-updated"
		expectedUserPtr := &expectedUser

		mockRepoUser.On("Update", mock.Anything, mock.MatchedBy(func(u *models_user.User) bool {
			return u.UID == inputUser.UID
		})).Return(expectedUserPtr, nil)

		result, err := service.Update(context.Background(), inputUser, nil)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserPtr, result)
		mockRepoUser.AssertExpectations(t)
	})

	t.Run("deve retornar erro para email inválido", func(t *testing.T) {
		_, _, service := setup()

		invalidUser := &models_user.User{
			Email: "invalid-email",
		}

		result, err := service.Update(context.Background(), invalidUser, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "email inválido")
	})

	t.Run("deve retornar erro de conflito de versão", func(t *testing.T) {
		mockRepoUser, _, service := setup()

		inputUser := &models_user.User{
			UID:     1,
			Email:   "valid@example.com",
			Version: 2, // versão enviada pelo front
		}

		// Simula retorno de conflito de versão pelo repositório
		mockRepoUser.On("Update", mock.Anything, inputUser).
			Return(nil, user_repositories.ErrVersionConflict).Once()

		result, err := service.Update(context.Background(), inputUser, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, user_repositories.ErrVersionConflict)

		mockRepoUser.AssertExpectations(t)
	})

	t.Run("deve lidar com usuário não encontrado", func(t *testing.T) {
		mockRepoUser, _, service := setup()

		user := &models_user.User{
			UID:     999,
			Email:   "valid@example.com",
			Version: 1,
		}

		mockRepoUser.On("Update", mock.Anything, mock.Anything).
			Return((*models_user.User)(nil), user_repositories.ErrUserNotFound)

		result, err := service.Update(context.Background(), user, nil)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, user_repositories.ErrUserNotFound))
		assert.Nil(t, result)

		mockRepoUser.AssertExpectations(t)
	})

	t.Run("deve lidar com outros erros do repositório", func(t *testing.T) {
		mockRepoUser, _, service := setup()

		user := &models_user.User{
			UID:     1,
			Email:   "valid@example.com",
			Version: 1,
		}

		mockRepoUser.On("Update", mock.Anything, mock.Anything).
			Return((*models_user.User)(nil), fmt.Errorf("erro no banco de dados"))

		result, err := service.Update(context.Background(), user, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar usuário")
		assert.Nil(t, result)

		mockRepoUser.AssertExpectations(t)
	})

	t.Run("endereço nil não tenta atualizar", func(t *testing.T) {
		mockRepoUser, mockAddressRepo, service := setup()

		user := &models_user.User{
			UID:     1,
			Email:   "valid@example.com",
			Version: 1,
		}

		mockRepoUser.On("Update", mock.Anything, user).
			Return(user, nil).Once()

		result, err := service.Update(context.Background(), user, nil)

		assert.NoError(t, err)
		assert.Equal(t, user, result)

		mockAddressRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		mockAddressRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		mockRepoUser.AssertExpectations(t)
	})

	t.Run("GetByID retorna ErrAddressNotFound", func(t *testing.T) {
		userRepo, addressRepo, service := setup()

		user := &models_user.User{UID: 1, Email: "teste@email.com", Version: 1}
		address := &models.Address{ID: 1}

		addressRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*models.Address)(nil), repositories_address.ErrAddressNotFound).Once()
		userRepo.On("Update", mock.Anything, user).
			Return(nil, nil).Once()

		result, err := service.Update(context.Background(), user, address)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, repositories_address.ErrAddressNotFound))

		addressRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("GetByID retorna erro genérico", func(t *testing.T) {
		userRepo, addressRepo, service := setup()

		user := &models_user.User{UID: 1, Email: "teste@email.com", Version: 1}
		address := &models.Address{ID: 1}
		erroGen := errors.New("erro qualquer")

		addressRepo.On("GetByID", mock.Anything, address.ID).
			Return((*models.Address)(nil), erroGen).Once()
		userRepo.On("Update", mock.Anything, user).
			Return(nil, nil).Once()

		result, err := service.Update(context.Background(), user, address)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "erro ao buscar o endereço")

		addressRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("Erro na atualização do endereço", func(t *testing.T) {
		mockRepo, addressRepo, service := setup()

		user := &models_user.User{UID: 1, Email: "test@example.com", Version: 1}
		address := &models.Address{ID: 1}
		existingAddr := &models.Address{ID: 1}

		addressRepo.On("GetByID", mock.Anything, address.ID).
			Return(existingAddr, nil).Once()

		addressRepo.On("Update", mock.Anything, address).
			Return(errors.New("falha update")).Once()

		mockRepo.On("Update", mock.Anything, mock.Anything).
			Return(user, nil).Once()

		result, err := service.Update(context.Background(), user, address)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "erro ao atualizar")

		addressRepo.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Atualização do endereço com sucesso", func(t *testing.T) {
		mockRepo, addressRepo, service := setup()

		user := &models_user.User{UID: 1, Email: "test@example.com", Version: 1}
		address := &models.Address{ID: 1}
		existingAddr := &models.Address{ID: 1}

		addressRepo.On("GetByID", mock.Anything, address.ID).
			Return(existingAddr, nil).Once()
		addressRepo.On("Update", mock.Anything, address).
			Return(nil).Once()

		mockRepo.On("Update", mock.Anything, mock.Anything).
			Return(user, nil).Once()

		result, err := service.Update(context.Background(), user, address)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		addressRepo.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {

	setup := func() (*user_repositories.MockUserRepository, UserService) {
		mockUserRepo := new(user_repositories.MockUserRepository)
		userService := NewUserService(
			mockUserRepo,
			new(addresses_repositories.MockAddressRepository),
			new(contact_repositories.MockContactRepository),
			new(user_category_relations_repositories.MockUserCategoryRelationRepo),
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
			Return(user_repositories.ErrUserNotFound)

		err := service.Delete(context.Background(), 999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), user_repositories.ErrUserNotFound.Error())
		assert.True(t, errors.Is(err, user_repositories.ErrUserNotFound), "deve envolver o erro original")
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
