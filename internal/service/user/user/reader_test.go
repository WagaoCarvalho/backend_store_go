package services

import (
	"context"
	"fmt"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	modelUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetAll(t *testing.T) {

	setup := func() (
		*mockUser.MockUser,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)

		userService := NewUserService(
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
		*mockUser.MockUser,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)

		userService := NewUserService(
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

func TestUserService_GetByEmail(t *testing.T) {

	setup := func() (
		*mockUser.MockUser,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)

		userService := NewUserService(
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
		*mockUser.MockUser,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)

		userService := NewUserService(
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
