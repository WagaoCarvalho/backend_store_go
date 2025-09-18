package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	modelsUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

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

	mockRepo := new(mockUser.MockUserRepository)
	mockHasher := new(mockHasher)
	mockToken := new(mockTokenGen)

	service := NewLoginService(mockRepo, mockToken, mockHasher)

	t.Run("sucesso", func(t *testing.T) {
		ctx := context.Background()
		email := "user@example.com"
		password := "123456"
		user := &modelsUser.User{UID: 1, Email: email, Password: "hashed", Status: true}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", password).Return(nil)
		mockToken.On("Generate", int64(1), email).Return("valid-token", nil)

		authResp, err := service.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.Equal(t, "valid-token", authResp.AccessToken)
		assert.Equal(t, "Bearer", authResp.TokenType)

		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("email inválido", func(t *testing.T) {
		authResp, err := service.Login(context.Background(), "invalid", "123")
		assert.Nil(t, authResp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
	})

	t.Run("usuário não encontrado", func(t *testing.T) {
		ctx := context.Background()
		email := "notfound@example.com"
		mockRepo.On("GetByEmail", ctx, email).Return((*modelsUser.User)(nil), errors.New("not found"))

		start := time.Now()
		authResp, err := service.Login(ctx, email, "123")
		elapsed := time.Since(start)

		assert.ErrorIs(t, err, errMsg.ErrCredentials)
		assert.Nil(t, authResp)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(1000)) // timing attack mitigation
		mockRepo.AssertExpectations(t)
	})

	t.Run("senha inválida", func(t *testing.T) {
		ctx := context.Background()
		email := "user@example.com"
		user := &modelsUser.User{UID: 1, Email: email, Password: "hashed", Status: true}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "wrong").Return(errors.New("wrong password"))

		authResp, err := service.Login(ctx, email, "wrong")

		assert.ErrorIs(t, err, errMsg.ErrCredentials)
		assert.Nil(t, authResp)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("conta desativada", func(t *testing.T) {
		ctx := context.Background()
		email := "inactive@example.com"
		user := &modelsUser.User{UID: 2, Email: email, Password: "hashed", Status: false}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "123").Return(nil)

		authResp, err := service.Login(ctx, email, "123")

		assert.ErrorIs(t, err, errMsg.ErrAccountDisabled)
		assert.Nil(t, authResp)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("erro ao gerar token", func(t *testing.T) {
		ctx := context.Background()
		email := "failtoken@example.com"
		user := &modelsUser.User{UID: 3, Email: email, Password: "hashed", Status: true}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "123").Return(nil)
		mockToken.On("Generate", int64(3), email).Return("", errors.New("gen error"))

		authResp, err := service.Login(ctx, email, "123")

		assert.ErrorIs(t, err, errMsg.ErrTokenGeneration)
		assert.Nil(t, authResp)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})
}
