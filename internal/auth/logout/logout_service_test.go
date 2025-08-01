package auth

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBlacklist struct {
	mock.Mock
}

func (m *mockBlacklist) Add(ctx context.Context, token string, duration time.Duration) error {
	args := m.Called(ctx, token, duration)
	return args.Error(0)
}

func (m *mockBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

type mockJWTService struct {
	mock.Mock
}

func (m *mockJWTService) Parse(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	token, _ := args.Get(0).(*jwt.Token)
	return token, args.Error(1)
}

func (m *mockJWTService) GetExpiration(token *jwt.Token) (time.Duration, error) {
	args := m.Called(token)
	return args.Get(0).(time.Duration), args.Error(1)
}

func generateTestToken(secret string, exp time.Time) string {
	claims := jwt.MapClaims{
		"sub":   "1",
		"email": "test@example.com",
		"exp":   exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestLogoutService_Logout(t *testing.T) {
	ctx := context.Background()
	log := logrus.New()
	log.SetOutput(io.Discard)
	loggerAdapter := logger.NewLoggerAdapter(log)
	secret := "test-secret"

	t.Run("logout com sucesso", func(t *testing.T) {
		tokenStr := generateTestToken(secret, time.Now().Add(time.Hour))
		mockBL := new(mockBlacklist)
		mockJWT := new(mockJWTService)

		token := &jwt.Token{Valid: true}
		mockJWT.On("Parse", tokenStr).Return(token, nil)
		mockJWT.On("GetExpiration", token).Return(1*time.Hour, nil)
		mockBL.On("Add", ctx, tokenStr, 1*time.Hour).Return(nil)

		service := NewLogoutService(mockBL, loggerAdapter, mockJWT)
		err := service.Logout(ctx, tokenStr)

		assert.NoError(t, err)
		mockBL.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("token malformado", func(t *testing.T) {
		mockBL := new(mockBlacklist)
		mockJWT := new(mockJWTService)

		mockJWT.On("Parse", "invalid-token").Return(nil, errors.New("malformed token"))

		service := NewLogoutService(mockBL, loggerAdapter, mockJWT)
		err := service.Logout(ctx, "invalid-token")

		assert.ErrorIs(t, err, ErrTokenValidation)
		mockJWT.AssertExpectations(t)
	})

	t.Run("token inválido", func(t *testing.T) {
		tokenStr := generateTestToken(secret, time.Now().Add(time.Hour))
		mockBL := new(mockBlacklist)
		mockJWT := new(mockJWTService)

		token := &jwt.Token{Valid: false}
		mockJWT.On("Parse", tokenStr).Return(token, nil)

		service := NewLogoutService(mockBL, loggerAdapter, mockJWT)
		err := service.Logout(ctx, tokenStr)

		assert.ErrorIs(t, err, ErrInvalidToken)
		mockJWT.AssertExpectations(t)
	})

	t.Run("erro ao pegar expiração", func(t *testing.T) {
		tokenStr := generateTestToken(secret, time.Now().Add(time.Hour))
		mockBL := new(mockBlacklist)
		mockJWT := new(mockJWTService)

		token := &jwt.Token{Valid: true}
		mockJWT.On("Parse", tokenStr).Return(token, nil)
		mockJWT.On("GetExpiration", token).Return(time.Duration(0), errors.New("claim exp faltando"))

		service := NewLogoutService(mockBL, loggerAdapter, mockJWT)
		err := service.Logout(ctx, tokenStr)

		assert.ErrorIs(t, err, ErrClaimExpInvalid)
		mockJWT.AssertExpectations(t)
	})

	t.Run("token expirado", func(t *testing.T) {
		tokenStr := generateTestToken(secret, time.Now().Add(-1*time.Hour))
		mockBL := new(mockBlacklist)
		mockJWT := new(mockJWTService)

		token := &jwt.Token{Valid: true}
		mockJWT.On("Parse", tokenStr).Return(token, nil)
		mockJWT.On("GetExpiration", token).Return(-1*time.Minute, nil)

		service := NewLogoutService(mockBL, loggerAdapter, mockJWT)
		err := service.Logout(ctx, tokenStr)

		assert.ErrorIs(t, err, ErrTokenExpired)
		mockJWT.AssertExpectations(t)
	})

	t.Run("erro ao adicionar na blacklist", func(t *testing.T) {
		tokenStr := generateTestToken(secret, time.Now().Add(time.Hour))
		mockBL := new(mockBlacklist)
		mockJWT := new(mockJWTService)

		token := &jwt.Token{Valid: true}
		mockJWT.On("Parse", tokenStr).Return(token, nil)
		mockJWT.On("GetExpiration", token).Return(1*time.Hour, nil)
		mockBL.On("Add", ctx, tokenStr, 1*time.Hour).Return(errors.New("falha redis"))

		service := NewLogoutService(mockBL, loggerAdapter, mockJWT)
		err := service.Logout(ctx, tokenStr)

		assert.ErrorIs(t, err, ErrBlacklistAdd)
		mockBL.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}
