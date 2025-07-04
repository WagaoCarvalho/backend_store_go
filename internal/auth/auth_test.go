package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/WagaoCarvalho/backend_store_go/internal/auth"
	modelsLogin "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	modelsUser "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
)

// Mock do TokenGenerator
type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) Generate(uid int64, email string) (string, error) {
	args := m.Called(uid, email)
	return args.String(0), args.Error(1)
}

// Mock do PasswordHasher
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Compare(hashed, plain string) error {
	args := m.Called(hashed, plain)
	return args.Error(0)
}

func TestLoginService_Login(t *testing.T) {
	mockRepo := new(repositories.MockUserRepository)
	mockJWT := new(MockJWTManager)
	mockHasher := new(MockHasher)

	service := auth.NewLoginService(mockRepo, mockJWT, mockHasher)

	ctx := context.Background()
	password := "senha123"
	hashedPassword := "hashed-senha123"

	user := &modelsUser.User{
		UID:      1,
		Email:    "teste@email.com",
		Password: hashedPassword,
		Status:   true,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, user.Email).Return(user, nil).Once()
		mockHasher.On("Compare", hashedPassword, password).Return(nil).Once()
		mockJWT.On("Generate", user.UID, user.Email).Return("token-123", nil).Once()

		token, err := service.Login(ctx, modelsLogin.LoginCredentials{
			Email:    user.Email,
			Password: password,
		})

		assert.NoError(t, err)
		assert.Equal(t, "token-123", token)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("invalid email format", func(t *testing.T) {
		token, err := service.Login(ctx, modelsLogin.LoginCredentials{
			Email:    "emailinvalido",
			Password: "senha123",
		})

		assert.ErrorIs(t, err, auth.ErrInvalidEmailFormat)
		assert.Empty(t, token)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, "naoexiste@email.com").
			Return(nil, repositories.ErrUserNotFound).Once()

		token, err := service.Login(ctx, modelsLogin.LoginCredentials{
			Email:    "naoexiste@email.com",
			Password: "senha123",
		})

		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, user.Email).Return(user, nil).Once()
		mockHasher.On("Compare", hashedPassword, "senhaErrada").Return(errors.New("senha incorreta")).Once()

		token, err := service.Login(ctx, modelsLogin.LoginCredentials{
			Email:    user.Email,
			Password: "senhaErrada",
		})

		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("account disabled", func(t *testing.T) {
		userDisabled := &modelsUser.User{
			UID:      2,
			Email:    "desativado@email.com",
			Password: hashedPassword,
			Status:   false,
		}

		mockRepo.On("GetByEmail", ctx, userDisabled.Email).Return(userDisabled, nil).Once()
		mockHasher.On("Compare", hashedPassword, password).Return(nil).Once()

		token, err := service.Login(ctx, modelsLogin.LoginCredentials{
			Email:    userDisabled.Email,
			Password: password,
		})

		assert.ErrorIs(t, err, auth.ErrAccountDisabled)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("error generating token", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, user.Email).Return(user, nil).Once()
		mockHasher.On("Compare", hashedPassword, password).Return(nil).Once()
		mockJWT.On("Generate", user.UID, user.Email).Return("", errors.New("falha gerar token")).Once()

		token, err := service.Login(ctx, modelsLogin.LoginCredentials{
			Email:    user.Email,
			Password: password,
		})

		assert.ErrorIs(t, err, auth.ErrTokenGeneration)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("unexpected error on GetByEmail", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, user.Email).Return(nil, errors.New("erro inesperado")).Once()

		token, err := service.Login(ctx, modelsLogin.LoginCredentials{
			Email:    user.Email,
			Password: password,
		})

		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}
