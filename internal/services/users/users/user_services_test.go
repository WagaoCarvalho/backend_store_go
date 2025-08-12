package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	model_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/users"
	repo_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/users"
	"github.com/WagaoCarvalho/backend_store_go/logger"
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
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
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

		user := &model_user.User{
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

		newUser := &model_user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "senha123",
			Status:   true,
		}

		hashed := "hashedSenha123"
		mockHasher.On("Hash", "senha123").Return(hashed, nil).Once()

		createdUser := &model_user.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashed,
			Status:   true,
		}

		mockUserRepo.
			On("Create", mock.Anything, mock.MatchedBy(func(u *model_user.User) bool {
				return u.Email == newUser.Email && u.Password == hashed
			})).
			Run(func(args mock.Arguments) {
				args.Get(1).(*model_user.User).UID = 1
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

		newUser := model_user.User{Email: "test@example.com"}
		mockUserRepo.On("Create", mock.Anything, &newUser).Return(nil, errors.New("erro no banco de dados"))

		_, err := userService.Create(context.Background(), &newUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar usuário")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("usuário criado é nulo", func(t *testing.T) {
		mockUserRepo, _, userService := setup()

		user := &model_user.User{
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

		_, err := userService.Create(context.Background(), &model_user.User{Email: "email-invalido"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
	})

}

func TestUserService_GetAll(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar todos os usuários com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUsers := []*model_user.User{
			{UID: 1, Username: "user1", Email: "user1@example.com", Status: true},
			{UID: 2, Username: "user2", Email: "user2@example.com", Status: false},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedUsers, nil)

		users, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao falhar no repositório", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetAll", mock.Anything).Return(nil, fmt.Errorf("erro ao acessar o banco"))

		users, err := service.GetAll(context.Background())

		assert.ErrorContains(t, err, "erro ao acessar o banco")
		assert.Nil(t, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar usuário quando encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUser := &model_user.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

		user, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro para ID inválido (<= 0)", func(t *testing.T) {
		mockRepo, _, service := setup()

		user, err := service.GetByID(context.Background(), 0) // ID inválido

		assert.Nil(t, user)
		assert.Error(t, err)
		assert.EqualError(t, err, "ID inválido")

		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("Deve retornar erro quando usuário não existe", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, fmt.Errorf("usuário não encontrado"))

		user, err := service.GetByID(context.Background(), 999)

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetVersionByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar versão quando usuário for encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil)

		version, err := service.GetVersionByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(999)).Return(
			int64(0),
			repo_user.ErrUserNotFound,
		)

		version, err := service.GetVersionByID(context.Background(), 999)

		assert.ErrorIs(t, err, repo_user.ErrUserNotFound)
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico quando falhar no repositório", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(2)).Return(
			int64(0),
			fmt.Errorf("erro no banco de dados"),
		)

		version, err := service.GetVersionByID(context.Background(), 2)

		assert.ErrorContains(t, err, "versão inválida")
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar usuário quando encontrado por e-mail", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUser := &model_user.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}

		mockRepo.On("GetByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

		user, err := service.GetByEmail(context.Background(), "user1@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando repositório falha ao buscar por e-mail", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetByEmail", mock.Anything, "inexistente@example.com").Return(
			nil,
			fmt.Errorf("usuário não encontrado"),
		)

		user, err := service.GetByEmail(context.Background(), "inexistente@example.com")

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByName(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar lista de usuários quando encontrados por nome parcial", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUsers := []*model_user.User{
			{
				UID:      1,
				Username: "user1",
				Email:    "user1@example.com",
				Status:   true,
			},
			{
				UID:      2,
				Username: "user123",
				Email:    "user123@example.com",
				Status:   true,
			},
		}

		mockRepo.On("GetByName", mock.Anything, "user").Return(expectedUsers, nil)

		users, err := service.GetByName(context.Background(), "user")

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando repositório falha ao buscar por nome", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetByName", mock.Anything, "inexistente").Return(
			nil,
			fmt.Errorf("usuário não encontrado"),
		)

		users, err := service.GetByName(context.Background(), "inexistente")

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Nil(t, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar erro ao atualizar com e-mail inválido", func(t *testing.T) {
		_, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "email_invalido",
			Version: 1,
		}

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Nil(t, updatedUser)
	})

	t.Run("Deve retornar erro ao atualizar com versão inválida", func(t *testing.T) {
		_, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 0,
		}

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, ErrInvalidVersion)
		assert.Nil(t, updatedUser)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, repo_user.ErrUserNotFound)

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, repo_user.ErrUserNotFound)
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de conflito de versão", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 2,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, repo_user.ErrVersionConflict)

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, repo_user.ErrVersionConflict)
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico ao atualizar", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, fmt.Errorf("erro interno"))

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorContains(t, err, "erro ao atualizar usuário")
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &model_user.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  1,
		}

		expected := &model_user.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  2,
		}

		mockRepo.On("Update", mock.Anything, user).Return(expected, nil)

		updatedUser, err := service.Update(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, expected, updatedUser)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Disable(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve desativar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Disable", mock.Anything, int64(1)).
			Return(nil).Once()

		err := service.Disable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao desativar usuário", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Disable", mock.Anything, int64(2)).
			Return(fmt.Errorf("falha no banco")).Once()

		err := service.Disable(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao desabilitar usuário")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Enable(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve ativar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Enable", mock.Anything, int64(1)).
			Return(nil).Once()

		err := service.Enable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando usuário não for encontrado ao habilitar", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Enable", mock.Anything, int64(42)).
			Return(repo.ErrUserNotFound).Once()

		err := service.Enable(context.Background(), 42)

		assert.ErrorIs(t, err, repo.ErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao ativar usuário", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Enable", mock.Anything, int64(2)).
			Return(fmt.Errorf("falha no banco")).Once()

		err := service.Enable(context.Background(), 2)

		assert.ErrorContains(t, err, "falha no banco")
		mockRepo.AssertExpectations(t)
	})

}

func TestUserService_Delete(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve deletar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao falhar na deleção", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(2)).Return(fmt.Errorf("erro no banco"))

		err := service.Delete(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao deletar usuário")
		mockRepo.AssertExpectations(t)
	})
}
