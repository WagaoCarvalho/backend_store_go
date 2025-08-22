package auth

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models_login "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/model/user"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models_user.User), args.Error(1)
}

// Mock do PasswordHasher
type MockHasher struct {
	mock.Mock
}

type mockHasher struct {
	mock.Mock
}

func (m *mockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockHasher) Compare(hashed, plain string) error {
	args := m.Called(hashed, plain)
	return args.Error(0)
}

type mockTokenGen struct{ mock.Mock }

func (m *mockTokenGen) Generate(uid int64, email string) (string, error) {
	args := m.Called(uid, email)
	return args.String(0), args.Error(1)
}

func TestLoginService_Login(t *testing.T) {
	log := logrus.New()
	log.SetOutput(io.Discard)
	adapter := logger.NewLoggerAdapter(log)

	mockRepo := new(repo.MockUserRepository)
	mockHasher := new(mockHasher)
	mockToken := new(mockTokenGen)

	service := NewLoginService(mockRepo, adapter, mockToken, mockHasher)

	t.Run("sucesso", func(t *testing.T) {
		ctx := context.Background()
		email := "user@example.com"
		password := "123456"
		user := &models_user.User{UID: 1, Email: email, Password: "hashed", Status: true}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", password).Return(nil)
		mockToken.On("Generate", int64(1), email).Return("valid-token", nil)

		token, err := service.Login(ctx, models_login.LoginCredentials{Email: email, Password: password})

		assert.NoError(t, err)
		assert.Equal(t, "valid-token", token)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("email inválido", func(t *testing.T) {
		token, err := service.Login(context.Background(), models_login.LoginCredentials{
			Email:    "invalid",
			Password: "123",
		})
		assert.ErrorIs(t, err, err_msg.ErrInvalidEmailFormat)
		assert.Empty(t, token)
	})

	t.Run("usuário não encontrado", func(t *testing.T) {
		ctx := context.Background()
		email := "notfound@example.com"
		mockRepo.On("GetByEmail", ctx, email).Return((*models_user.User)(nil), errors.New("not found"))

		start := time.Now()
		token, err := service.Login(ctx, models_login.LoginCredentials{Email: email, Password: "123"})
		elapsed := time.Since(start)

		assert.ErrorIs(t, err, err_msg.ErrInvalidCredentials)
		assert.Empty(t, token)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(1000))
		mockRepo.AssertExpectations(t)
	})

	t.Run("senha inválida", func(t *testing.T) {
		ctx := context.Background()
		email := "user@example.com"
		user := &models_user.User{UID: 1, Email: email, Password: "hashed", Status: true}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "wrong").Return(errors.New("wrong password"))

		token, err := service.Login(ctx, models_login.LoginCredentials{Email: email, Password: "wrong"})

		assert.ErrorIs(t, err, err_msg.ErrInvalidCredentials)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("conta desativada", func(t *testing.T) {
		ctx := context.Background()
		email := "inactive@example.com"
		user := &models_user.User{UID: 2, Email: email, Password: "hashed", Status: false}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "123").Return(nil)

		token, err := service.Login(ctx, models_login.LoginCredentials{Email: email, Password: "123"})

		assert.ErrorIs(t, err, err_msg.ErrAccountDisabled)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("erro ao gerar token", func(t *testing.T) {
		ctx := context.Background()
		email := "failtoken@example.com"
		user := &models_user.User{UID: 3, Email: email, Password: "hashed", Status: true}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "123").Return(nil)
		mockToken.On("Generate", int64(3), email).Return("", errors.New("gen error"))

		token, err := service.Login(ctx, models_login.LoginCredentials{Email: email, Password: "123"})

		assert.ErrorIs(t, err, err_msg.ErrTokenGeneration)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})
}
