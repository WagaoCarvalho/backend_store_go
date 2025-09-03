package auth

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mock_user "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/login"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
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

func TestLoginService_LoginDTO(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	adapter := logger.NewLoggerAdapter(log)

	mockRepo := new(mock_user.MockUserRepository)
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

		credentialsDTO := dto.LoginCredentialsDTO{Email: email, Password: password}
		authRespDTO, err := service.Login(ctx, credentialsDTO)

		assert.NoError(t, err)
		assert.Equal(t, "valid-token", authRespDTO.AccessToken)
		assert.Equal(t, "Bearer", authRespDTO.TokenType)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("email inválido", func(t *testing.T) {
		credentialsDTO := dto.LoginCredentialsDTO{Email: "invalid", Password: "123"}
		authRespDTO, err := service.Login(context.Background(), credentialsDTO)
		assert.ErrorIs(t, err, err_msg.ErrEmailFormat)
		assert.Nil(t, authRespDTO)
	})

	t.Run("usuário não encontrado", func(t *testing.T) {
		ctx := context.Background()
		email := "notfound@example.com"
		mockRepo.On("GetByEmail", ctx, email).Return((*models_user.User)(nil), errors.New("not found"))

		start := time.Now()
		credentialsDTO := dto.LoginCredentialsDTO{Email: email, Password: "123"}
		authRespDTO, err := service.Login(ctx, credentialsDTO)
		elapsed := time.Since(start)

		assert.ErrorIs(t, err, err_msg.ErrCredentials)
		assert.Nil(t, authRespDTO)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(1000))
		mockRepo.AssertExpectations(t)
	})

	t.Run("senha inválida", func(t *testing.T) {
		ctx := context.Background()
		email := "user@example.com"
		user := &models_user.User{UID: 1, Email: email, Password: "hashed", Status: true}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "wrong").Return(errors.New("wrong password"))

		credentialsDTO := dto.LoginCredentialsDTO{Email: email, Password: "wrong"}
		authRespDTO, err := service.Login(ctx, credentialsDTO)

		assert.ErrorIs(t, err, err_msg.ErrCredentials)
		assert.Nil(t, authRespDTO)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("conta desativada", func(t *testing.T) {
		ctx := context.Background()
		email := "inactive@example.com"
		user := &models_user.User{UID: 2, Email: email, Password: "hashed", Status: false}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)
		mockHasher.On("Compare", "hashed", "123").Return(nil)

		credentialsDTO := dto.LoginCredentialsDTO{Email: email, Password: "123"}
		authRespDTO, err := service.Login(ctx, credentialsDTO)

		assert.ErrorIs(t, err, err_msg.ErrAccountDisabled)
		assert.Nil(t, authRespDTO)
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

		credentialsDTO := dto.LoginCredentialsDTO{Email: email, Password: "123"}
		authRespDTO, err := service.Login(ctx, credentialsDTO)

		assert.ErrorIs(t, err, err_msg.ErrTokenGeneration)
		assert.Nil(t, authRespDTO)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})
}
