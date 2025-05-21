package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	models_login "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestLoginService_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(services.MockUserRepository)
		service := NewLoginService(mockRepo)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

		user := models.User{
			UID:      1,
			Email:    "teste@email.com",
			Password: string(hashedPassword),
			Status:   true,
		}

		mockRepo.On("GetByEmail", mock.Anything, "teste@email.com").Return(user, nil)

		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "teste@email.com",
			Password: "senha123",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid email format", func(t *testing.T) {
		mockRepo := new(services.MockUserRepository)
		service := NewLoginService(mockRepo)

		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "email-invalido",
			Password: "senha123",
		})

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "formato de email inválido")
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo := new(services.MockUserRepository)
		service := NewLoginService(mockRepo)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("outrasenha"), bcrypt.DefaultCost)

		user := models.User{
			UID:      1,
			Email:    "teste@email.com",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByEmail", mock.Anything, "teste@email.com").Return(user, nil)

		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "teste@email.com",
			Password: "senhaErrada",
		})

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "credenciais inválidas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(services.MockUserRepository)
		service := NewLoginService(mockRepo)

		mockRepo.On("GetByEmail", mock.Anything, "naoexiste@email.com").
			Return(models.User{}, repositories.ErrUserNotFound)

		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "naoexiste@email.com",
			Password: "senha123",
		})

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "credenciais inválidas")
		mockRepo.AssertExpectations(t)
	})

	t.Run("user disabled", func(t *testing.T) {
		mockRepo := new(services.MockUserRepository)
		service := NewLoginService(mockRepo)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

		user := models.User{
			UID:      1,
			Email:    "teste@email.com",
			Password: string(hashedPassword),
			Status:   false,
		}

		mockRepo.On("GetByEmail", mock.Anything, "teste@email.com").Return(user, nil)

		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "teste@email.com",
			Password: "senha123",
		})

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "conta desativada")
		mockRepo.AssertExpectations(t)
	})

	t.Run("unexpected error on GetByEmail", func(t *testing.T) {
		mockRepo := new(services.MockUserRepository)
		service := NewLoginService(mockRepo)

		mockRepo.On("GetByEmail", mock.Anything, "teste@email.com").
			Return(models.User{}, errors.New("falha no banco"))

		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "teste@email.com",
			Password: "senha123",
		})

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "erro ao buscar usuário")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error generating token", func(t *testing.T) {
		mockRepo := new(services.MockUserRepository)

		// Função mockada de geração de token
		fakeJWTGenerator := func(uid int64, email string) (string, error) {
			return "", fmt.Errorf("falha ao gerar token")
		}

		service := NewLoginServiceWithJWT(mockRepo, fakeJWTGenerator)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

		user := models.User{
			UID:      1,
			Email:    "teste@email.com",
			Password: string(hashedPassword),
			Status:   true,
		}

		mockRepo.On("GetByEmail", mock.Anything, "teste@email.com").Return(user, nil)

		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "teste@email.com",
			Password: "senha123",
		})

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "erro ao gerar token de acesso")
		mockRepo.AssertExpectations(t)
	})
}
