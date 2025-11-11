package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
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
	setup := func() (*mockUser.MockUser, *MockHasher, User) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)
		userService := NewUserService(mockUserRepo, mockHasher)
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

func TestUserService_Update(t *testing.T) {
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

	t.Run("Deve retornar erro ao atualizar com e-mail inválido", func(t *testing.T) {
		_, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "email_invalido",
			Version: 1,
		}

		err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("Deve retornar erro ao atualizar com versão inválida", func(t *testing.T) {
		_, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 0,
		}

		err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(errMsg.ErrNotFound)

		err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de conflito de versão", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 2,
		}

		mockRepo.On("Update", mock.Anything, user).Return(errMsg.ErrVersionConflict)

		err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico ao atualizar", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(fmt.Errorf("erro interno"))

		err := service.Update(context.Background(), user)

		assert.ErrorContains(t, err, "erro ao atualizar")
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

		mockRepo.On("Update", mock.Anything, user).Return(nil)

		err := service.Update(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {

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
