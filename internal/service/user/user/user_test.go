package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	modelUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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
	return nil
}

func TestUserService_Create(t *testing.T) {
	setup := func() (*mockUser.MockUserRepository, *MockHasher, User) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)
		userService := NewUser(mockUserRepo, mockHasher)
		return mockUserRepo, mockHasher, userService
	}

	t.Run("erro na validação do usuário", func(t *testing.T) {
		_, _, userService := setup()

		user := &modelUser.User{
			Email:    "", // inválido
			Username: "", // inválido
			Password: "", // inválido
		}

		result, err := userService.Create(context.Background(), user)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("erro ao hashear senha", func(t *testing.T) {
		mockRepo, mockHasher, userService := setup()

		user := &modelUser.User{
			Email:    "teste@example.com",
			Username: "teste",
			Password: "Senha@123",
			Status:   true,
		}

		// Simula falha no hash
		mockHasher.On("Hash", "Senha@123").Return("", errors.New("falha no hash")).Once()

		result, err := userService.Create(context.Background(), user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao hashear senha")

		mockHasher.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create") // Create não deve ser chamado
	})

	t.Run("sucesso ao criar usuário", func(t *testing.T) {
		mockRepo, mockHasher, userService := setup()

		user := &modelUser.User{
			Email:    "teste@example.com",
			Username: "teste",
			Password: "Senha@123",
			Status:   true, // necessário para passar validação
		}

		// Hash da senha
		hashed := "hashedSenha123"
		mockHasher.On("Hash", "Senha@123").Return(hashed, nil).Once()

		// Mock do Create retornando UID preenchido
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *modelUser.User) bool {
			return u.Email == user.Email && u.Password == hashed
		})).
			Return(&modelUser.User{UID: 1, Email: user.Email, Password: hashed}, nil).Once()

		result, err := userService.Create(context.Background(), user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.UID)

		mockHasher.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao criar usuário no repo", func(t *testing.T) {
		mockRepo, mockHasher, userService := setup()

		user := &modelUser.User{
			Email:    "teste@example.com",
			Username: "teste",
			Password: "Senha@123",
			Status:   true,
		}

		mockHasher.On("Hash", "Senha@123").Return("hashedSenha123", nil).Once()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("erro no banco")).Once()

		result, err := userService.Create(context.Background(), user)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar")

		mockHasher.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("usuário criado é nulo", func(t *testing.T) {
		mockRepo, mockHasher, userService := setup()

		user := &modelUser.User{
			Email:    "teste@example.com",
			Username: "teste",
			Password: "Senha@123",
			Status:   true,
		}

		mockHasher.On("Hash", "Senha@123").Return("hashedSenha123", nil).Once()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, nil).Once()

		result, err := userService.Create(context.Background(), user)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "usuário criado é nulo")

		mockHasher.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAll(t *testing.T) {

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,
			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar todos os usuários com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUsers := []*modelUser.User{
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

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar usuário quando encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUser := &modelUser.User{
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
		assert.Equal(t, err, errMsg.ErrZeroID)

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

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		_, _, service := setup()

		version, err := service.GetVersionByID(context.Background(), 0)

		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

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
			errMsg.ErrNotFound,
		)

		version, err := service.GetVersionByID(context.Background(), 999)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
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

		assert.ErrorContains(t, err, "conflito de versão")
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: email inválido", func(t *testing.T) {
		_, _, service := setup()

		user, err := service.GetByEmail(context.Background(), "   ")

		assert.Nil(t, user)
		assert.EqualError(t, err, "email inválido")
	})

	t.Run("Deve retornar usuário quando encontrado por e-mail", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUser := &modelUser.User{
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

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: nome inválido", func(t *testing.T) {
		_, _, service := setup()

		user, err := service.GetByName(context.Background(), "   ")

		assert.Nil(t, user)
		assert.EqualError(t, err, "nome inválido")
	})

	t.Run("Deve retornar lista de usuários quando encontrados por nome parcial", func(t *testing.T) {
		mockRepo, _, service := setup()

		expectedUsers := []*modelUser.User{
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

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("Deve retornar erro ao atualizar com e-mail inválido", func(t *testing.T) {
		_, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "email_invalido",
			Version: 1,
		}

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Nil(t, updatedUser)
	})

	t.Run("Deve retornar erro ao atualizar com versão inválida", func(t *testing.T) {
		_, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 0,
		}

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		assert.Nil(t, updatedUser)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, errMsg.ErrNotFound)

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de conflito de versão", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 2,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, errMsg.ErrVersionConflict)

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico ao atualizar", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, fmt.Errorf("erro interno"))

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorContains(t, err, "erro ao atualizar")
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  1,
		}

		expected := &modelUser.User{
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

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		_, _, service := setup()

		err := service.Disable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

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
			Return(fmt.Errorf("erro ao desabilitar")).Once()

		err := service.Disable(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao desabilitar")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Enable(t *testing.T) {

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		_, _, service := setup()

		err := service.Enable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

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
			Return(errMsg.ErrNotFound).Once()

		err := service.Enable(context.Background(), 42)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
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

	setup := func() (
		*mockUser.MockUserRepository,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUserRepository)
		mockHasher := new(MockHasher)

		userService := NewUser(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		_, _, service := setup()

		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("Deve deletar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao falhar na deleção", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(2)).Return(fmt.Errorf("erro ao deletar"))

		err := service.Delete(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao deletar")
		mockRepo.AssertExpectations(t)
	})
}
