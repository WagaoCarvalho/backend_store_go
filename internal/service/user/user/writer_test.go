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
			Email:    "",
			Username: "",
			Password: "",
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

		mockHasher.On("Hash", "Senha@123").Return("", errors.New("falha no hash")).Once()

		result, err := userService.Create(context.Background(), user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao hashear senha")

		mockHasher.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("sucesso ao criar usuário", func(t *testing.T) {
		mockRepo, mockHasher, userService := setup()

		user := &modelUser.User{
			Email:    "teste@example.com",
			Username: "teste",
			Password: "Senha@123",
			Status:   true,
		}

		hashed := "hashedSenha123"
		mockHasher.On("Hash", "Senha@123").Return(hashed, nil).Once()

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
		mockRepo, mockHasher, service := setup()
		user := &modelUser.User{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "SecurePass123!", // Senha com complexidade correta
			Description: "Test user",
			Status:      true,
		}

		// Mock do hasher
		hashedPassword := "$2a$10$hashedpassword"
		mockHasher.On("Hash", user.Password).Return(hashedPassword, nil)

		// Mock do repo retornando erro
		repoError := fmt.Errorf("erro no banco")
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *modelUser.User) bool {
			return u.Password == hashedPassword
		})).Return(nil, repoError)

		result, err := service.Create(context.Background(), user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), "erro no banco")
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

	t.Run("Deve retornar erro quando usuário é nil", func(t *testing.T) {
		_, _, service := setup()
		err := service.Update(context.Background(), nil)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "usuário não pode ser nulo")
	})

	t.Run("Deve retornar erro quando ID é zero ou negativo", func(t *testing.T) {
		_, _, service := setup()
		user := &modelUser.User{
			UID:      0,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  1,
		}
		err := service.Update(context.Background(), user)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		assert.Contains(t, err.Error(), "ID do usuário inválido")
	})

	t.Run("Deve retornar erro quando versão é negativa", func(t *testing.T) {
		_, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  -1,
		}
		err := service.Update(context.Background(), user)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "versão não pode ser negativa")
	})

	t.Run("Deve retornar erro ao atualizar com e-mail inválido", func(t *testing.T) {
		_, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "email_invalido",
			Version:  1,
		}
		err := service.Update(context.Background(), user)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "email inválido")
	})

	t.Run("Deve retornar erro ao atualizar com username muito curto", func(t *testing.T) {
		_, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "ab", // Menor que 3 caracteres
			Email:    "user@example.com",
			Version:  1,
		}
		err := service.Update(context.Background(), user)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "deve ter entre 3 e 50 caracteres")
	})

	t.Run("Deve hash da senha quando fornecida e não for hash", func(t *testing.T) {
		mockRepo, mockHasher, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Password: "SecurePass123!", // Senha com complexidade
			Version:  1,
		}

		// Configurar o mock do hasher
		expectedHash := "$2a$10$hashedpassword"
		mockHasher.On("Hash", "SecurePass123!").Return(expectedHash, nil)
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *modelUser.User) bool {
			return u.Password == expectedHash
		})).Return(nil)

		err := service.Update(context.Background(), user)
		assert.NoError(t, err)
		mockHasher.AssertCalled(t, "Hash", "SecurePass123!")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Não deve hash da senha quando já for hash", func(t *testing.T) {
		mockRepo, mockHasher, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Password: "$2a$10$alreadyhashed",
			Version:  1,
		}

		// O hasher NÃO deve ser chamado
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *modelUser.User) bool {
			return u.Password == "$2a$10$alreadyhashed"
		})).Return(nil)

		err := service.Update(context.Background(), user)
		assert.NoError(t, err)
		mockHasher.AssertNotCalled(t, "Hash", mock.Anything)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Deve retornar erro interno ao falhar o hash da senha", func(t *testing.T) {
		mockRepo, mockHasher, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Password: "SecurePass123!",
			Version:  1,
		}

		mockHasher.On("Hash", "SecurePass123!").Return("", fmt.Errorf("erro no hash"))

		err := service.Update(context.Background(), user)
		assert.ErrorIs(t, err, errMsg.ErrInternal)
		assert.Contains(t, err.Error(), "erro ao processar senha")
		mockHasher.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(errMsg.ErrNotFound)
		err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Contains(t, err.Error(), "usuário não encontrado")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de conflito de versão", func(t *testing.T) {
		mockRepo, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(errMsg.ErrVersionConflict)
		err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		assert.Contains(t, err.Error(), "versão conflitante")
		assert.Contains(t, err.Error(), "dados desatualizados")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico ao atualizar", func(t *testing.T) {
		mockRepo, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  1,
		}

		repoError := fmt.Errorf("erro no banco de dados")
		mockRepo.On("Update", mock.Anything, user).Return(repoError)
		err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "erro no banco de dados")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com sucesso sem alterar senha", func(t *testing.T) {
		mockRepo, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Password: "", // Senha vazia
			Version:  1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil)
		err := service.Update(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com sucesso com senha vazia mas diferente de hash", func(t *testing.T) {
		mockRepo, mockHasher, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Password: "", // Senha vazia - não deve validar complexidade
			Version:  1,
		}

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *modelUser.User) bool {
			return u.Password == ""
		})).Return(nil)

		err := service.Update(context.Background(), user)
		assert.NoError(t, err)
		mockHasher.AssertNotCalled(t, "Hash", mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com descrição vazia", func(t *testing.T) {
		mockRepo, _, service := setup()
		user := &modelUser.User{
			UID:         1,
			Username:    "usuario",
			Email:       "user@example.com",
			Description: "",
			Version:     1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil)
		err := service.Update(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com status false", func(t *testing.T) {
		mockRepo, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Status:   false,
			Version:  1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil)
		err := service.Update(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com versão zero (para novos registros)", func(t *testing.T) {
		mockRepo, _, service := setup()
		user := &modelUser.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  0, // Versão zero é válida (para novos registros ou reset)
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
