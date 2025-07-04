package services

import (
	"context"
	"errors"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	user_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Compare(_, _ string) error {
	// Implementado apenas para satisfazer a interface
	return nil
}

func TestUserService_Create(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New()) // logger real

	setup := func() (
		*user_repositories.MockUserRepository,
		*MockHasher,
		*userService,
	) {
		mockUserRepo := new(user_repositories.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("erro ao hashear senha", func(t *testing.T) {
		mockUserRepo, mockHasher, userService := setup()

		user := &models_user.User{
			Email:    "test@example.com",
			Password: "senhaInvalidaParaHash",
		}

		mockHasher.On("Hash", "senhaInvalidaParaHash").Return("", errors.New("falha no hash")).Once()

		_, err := userService.Create(context.Background(), user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao hashear senha")
		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("sucesso ao criar usuário com todos os dados", func(t *testing.T) {
		mockUserRepo, mockHasher, userService := setup()

		newUser := &models_user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "senha123",
			Status:   true,
		}

		hashed := "hashedSenha123"
		mockHasher.On("Hash", "senha123").Return(hashed, nil).Once()

		createdUser := &models_user.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashed,
			Status:   true,
		}

		mockUserRepo.
			On("Create", mock.Anything, mock.MatchedBy(func(u *models_user.User) bool {
				return u.Email == newUser.Email && u.Password == hashed
			})).
			Run(func(args mock.Arguments) {
				args.Get(1).(*models_user.User).UID = 1
			}).
			Return(createdUser, nil)

		result, err := userService.Create(context.Background(), newUser)

		assert.NoError(t, err)
		assert.Equal(t, createdUser, result)
		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
	t.Run("erro ao criar usuário", func(t *testing.T) {
		mockUserRepo, _, userService := setup()

		newUser := models_user.User{Email: "test@example.com"}
		mockUserRepo.On("Create", mock.Anything, &newUser).Return(nil, errors.New("erro no banco de dados"))

		_, err := userService.Create(context.Background(), &newUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar usuário")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("usuário criado é nulo", func(t *testing.T) {
		mockUserRepo, _, userService := setup()

		user := &models_user.User{
			Email: "valid@email.com",
		}

		mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil, nil)

		_, err := userService.Create(context.Background(), user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "usuário criado é nulo")
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("email inválido", func(t *testing.T) {
		_, _, userService := setup()

		_, err := userService.Create(context.Background(), &models_user.User{Email: "email-invalido"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
	})

}

// func TestUserService_GetUsers(t *testing.T) {
// 	mockUserRepo := new(user_repositories.MockUserRepository)

// 	expectedUsers := []*models_user.User{
// 		{UID: 1, Username: "user1", Email: "user1@example.com", Status: true},
// 		{UID: 2, Username: "user2", Email: "user2@example.com", Status: false},
// 	}
// 	mockUserRepo.On("GetAll", mock.Anything).Return(expectedUsers, nil)

// 	userService := NewUserService(mockUserRepo)
// 	users, err := userService.GetAll(context.Background())

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedUsers, users)
// 	mockUserRepo.AssertExpectations(t)
// }

// func TestUserService_GetUserByID(t *testing.T) {

// 	setupMocks := func() *user_repositories.MockUserRepository {
// 		return new(user_repositories.MockUserRepository)
// 	}

// 	t.Run("Deve retornar usuário quando encontrado", func(t *testing.T) {
// 		mockUserRepo := setupMocks()

// 		expectedUser := &models_user.User{
// 			UID:      1,
// 			Username: "user1",
// 			Email:    "user1@example.com",
// 			Status:   true,
// 		}

// 		mockUserRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

// 		userService := NewUserService(mockUserRepo)
// 		user, err := userService.GetByID(context.Background(), 1)

// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedUser, user)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("Deve retornar erro quando usuário não existe", func(t *testing.T) {
// 		mockUserRepo := setupMocks()

// 		mockUserRepo.On("GetByID", mock.Anything, int64(999)).Return(
// 			nil, // ponteiro nil
// 			fmt.Errorf("usuário não encontrado"),
// 		)

// 		userService := NewUserService(mockUserRepo)
// 		user, err := userService.GetByID(context.Background(), 999)

// 		assert.ErrorContains(t, err, "usuário não encontrado")
// 		assert.Nil(t, user) // agora user deve ser nil
// 		mockUserRepo.AssertExpectations(t)
// 	})

// }

// func TestUserService_GetVersionByID(t *testing.T) {
// 	t.Run("deve retornar a versão corretamente", func(t *testing.T) {
// 		mockRepo := new(user_repositories.MockUserRepository)
// 		service := NewUserService(
// 			mockRepo)

// 		uid := int64(1)
// 		expectedVersion := int64(5)

// 		mockRepo.On("GetVersionByID", mock.Anything, uid).Return(expectedVersion, nil).Once()

// 		version, err := service.GetVersionByID(context.Background(), uid)

// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedVersion, version)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("deve retornar erro de usuário não encontrado", func(t *testing.T) {
// 		mockRepo := new(user_repositories.MockUserRepository)
// 		service := NewUserService(
// 			mockRepo)

// 		uid := int64(999)

// 		mockRepo.On("GetVersionByID", mock.Anything, uid).
// 			Return(int64(0), user_repositories.ErrUserNotFound).Once()

// 		version, err := service.GetVersionByID(context.Background(), uid)

// 		assert.ErrorIs(t, err, user_repositories.ErrUserNotFound)
// 		assert.Equal(t, int64(0), version)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("deve retornar erro genérico do repositório", func(t *testing.T) {
// 		mockRepo := new(user_repositories.MockUserRepository)
// 		service := NewUserService(
// 			mockRepo)

// 		uid := int64(2)
// 		repoErr := errors.New("falha no banco")

// 		mockRepo.On("GetVersionByID", mock.Anything, uid).Return(int64(0), repoErr).Once()

// 		version, err := service.GetVersionByID(context.Background(), uid)

// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "user: erro ao obter versão")
// 		assert.Equal(t, int64(0), version)
// 		mockRepo.AssertExpectations(t)
// 	})
// }

// func TestUserService_GetUserByEmail(t *testing.T) {

// 	setup := func() (*user_repositories.MockUserRepository, UserService) {
// 		mockUserRepo := new(user_repositories.MockUserRepository)

// 		service := NewUserService(
// 			mockUserRepo)
// 		return mockUserRepo, service
// 	}

// 	t.Run("Deve retornar usuário quando email existe", func(t *testing.T) {
// 		mockUserRepo, userService := setup()

// 		expectedUser := &models_user.User{
// 			UID:      1,
// 			Username: "user1",
// 			Email:    "user1@example.com",
// 			Status:   true,
// 		}

// 		mockUserRepo.On("GetByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

// 		user, err := userService.GetByEmail(context.Background(), "user1@example.com")

// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedUser, user)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("Deve retornar erro quando email não existe", func(t *testing.T) {
// 		mockUserRepo, userService := setup()

// 		mockUserRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(
// 			nil,
// 			fmt.Errorf("usuário não encontrado"),
// 		)

// 		user, err := userService.GetByEmail(context.Background(), "notfound@example.com")

// 		assert.ErrorContains(t, err, "usuário não encontrado")
// 		assert.Nil(t, user)
// 		mockUserRepo.AssertExpectations(t)
// 	})
// }

// func TestUserService_Update(t *testing.T) {
// 	setup := func() (*user_repositories.MockUserRepository, UserService) {
// 		mockUserRepo := new(user_repositories.MockUserRepository)

// 		service := NewUserService(
// 			mockUserRepo)

// 		return mockUserRepo, service
// 	}

// 	t.Run("versão inválida", func(t *testing.T) {
// 		_, service := setup()

// 		user := &models_user.User{
// 			UID:      1,
// 			Username: "user1",
// 			Email:    "valid@example.com",
// 			Status:   true,
// 		}

// 		updated, err := service.Update(context.Background(), user)

// 		assert.Nil(t, updated)
// 		assert.ErrorIs(t, err, ErrInvalidVersion)
// 	})

// 	t.Run("deve atualizar usuário com sucesso", func(t *testing.T) {
// 		mockRepoUser, service := setup()

// 		inputUser := &models_user.User{
// 			UID:      1,
// 			Username: "user1",
// 			Email:    "valid@example.com",
// 			Status:   true,
// 			Version:  1,
// 		}

// 		expectedUser := *inputUser
// 		expectedUser.Username = "user1-updated"
// 		expectedUserPtr := &expectedUser

// 		mockRepoUser.On("Update", mock.Anything, mock.MatchedBy(func(u *models_user.User) bool {
// 			return u.UID == inputUser.UID
// 		})).Return(expectedUserPtr, nil)

// 		result, err := service.Update(context.Background(), inputUser)

// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedUserPtr, result)
// 		mockRepoUser.AssertExpectations(t)
// 	})

// 	t.Run("deve retornar erro para email inválido", func(t *testing.T) {
// 		_, service := setup()

// 		invalidUser := &models_user.User{
// 			Email:   "invalid-email",
// 			Version: 1,
// 		}

// 		result, err := service.Update(context.Background(), invalidUser)

// 		assert.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Contains(t, err.Error(), "email inválido")
// 	})

// 	t.Run("deve retornar erro de conflito de versão", func(t *testing.T) {
// 		mockRepoUser, service := setup()

// 		inputUser := &models_user.User{
// 			UID:     1,
// 			Email:   "valid@example.com",
// 			Version: 2,
// 		}

// 		mockRepoUser.On("Update", mock.Anything, inputUser).
// 			Return(nil, user_repositories.ErrVersionConflict).Once()

// 		result, err := service.Update(context.Background(), inputUser)

// 		assert.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.ErrorIs(t, err, user_repositories.ErrVersionConflict)

// 		mockRepoUser.AssertExpectations(t)
// 	})

// 	t.Run("deve lidar com usuário não encontrado", func(t *testing.T) {
// 		mockRepoUser, service := setup()

// 		user := &models_user.User{
// 			UID:     999,
// 			Email:   "valid@example.com",
// 			Version: 1,
// 		}

// 		mockRepoUser.On("Update", mock.Anything, mock.Anything).
// 			Return((*models_user.User)(nil), user_repositories.ErrUserNotFound)

// 		result, err := service.Update(context.Background(), user)

// 		assert.Error(t, err)
// 		assert.True(t, errors.Is(err, user_repositories.ErrUserNotFound))
// 		assert.Nil(t, result)

// 		mockRepoUser.AssertExpectations(t)
// 	})

// 	t.Run("deve lidar com outros erros do repositório", func(t *testing.T) {
// 		mockRepoUser, service := setup()

// 		user := &models_user.User{
// 			UID:     1,
// 			Email:   "valid@example.com",
// 			Version: 1,
// 		}

// 		mockRepoUser.On("Update", mock.Anything, mock.Anything).
// 			Return((*models_user.User)(nil), fmt.Errorf("erro no banco de dados"))

// 		result, err := service.Update(context.Background(), user)

// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "erro ao atualizar usuário")
// 		assert.Nil(t, result)

// 		mockRepoUser.AssertExpectations(t)
// 	})
// }

// func TestUserService_Delete(t *testing.T) {

// 	setup := func() (*user_repositories.MockUserRepository, UserService) {
// 		mockUserRepo := new(user_repositories.MockUserRepository)
// 		userService := NewUserService(
// 			mockUserRepo)
// 		return mockUserRepo, userService
// 	}

// 	t.Run("deve deletar usuário com sucesso", func(t *testing.T) {
// 		mockUserRepo, service := setup()

// 		mockUserRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

// 		err := service.Delete(context.Background(), 1)

// 		assert.NoError(t, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("deve retornar erro quando usuário não existe", func(t *testing.T) {
// 		mockUserRepo, service := setup()

// 		mockUserRepo.On("Delete", mock.Anything, int64(999)).
// 			Return(user_repositories.ErrUserNotFound)

// 		err := service.Delete(context.Background(), 999)

// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), user_repositories.ErrUserNotFound.Error())
// 		assert.True(t, errors.Is(err, user_repositories.ErrUserNotFound), "deve envolver o erro original")
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("deve retornar erro genérico do repositório", func(t *testing.T) {
// 		mockUserRepo, service := setup()

// 		expectedErr := fmt.Errorf("erro no banco de dados")
// 		mockUserRepo.On("Delete", mock.Anything, int64(1)).
// 			Return(expectedErr)

// 		err := service.Delete(context.Background(), 1)

// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "erro ao deletar usuário")
// 		assert.True(t, errors.Is(err, expectedErr), "deve envolver o erro original")
// 		mockUserRepo.AssertExpectations(t)
// 	})
// }
