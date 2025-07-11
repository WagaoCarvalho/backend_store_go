package auth

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type customClaims struct{}

func (c customClaims) Valid() error {
	return nil
}

func (c customClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (c customClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c customClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c customClaims) GetIssuer() (string, error) {
	return "", nil
}

func (c customClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c customClaims) GetSubject() (string, error) {
	return "", nil
}

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

func generateTestToken(secret string, exp time.Time) string {
	claims := jwt.MapClaims{
		"uid":   1,
		"email": "test@example.com",
		"exp":   exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestLogoutService_Logout(t *testing.T) {
	log := logrus.New()
	log.SetOutput(io.Discard)
	loggerAdapter := logger.NewLoggerAdapter(log)
	secretKey := "test-secret"
	ctx := context.Background()

	t.Run("logout com sucesso", func(t *testing.T) {
		token := generateTestToken(secretKey, time.Now().Add(time.Hour))
		mockBL := new(mockBlacklist)
		mockBL.On("Add", mock.Anything, token, mock.AnythingOfType("time.Duration")).Return(nil)

		service := NewLogoutService(mockBL, loggerAdapter, secretKey)
		err := service.Logout(ctx, token)

		assert.NoError(t, err)
		mockBL.AssertExpectations(t)
	})

	t.Run("token malformado", func(t *testing.T) {
		mockBL := new(mockBlacklist)
		service := NewLogoutService(mockBL, loggerAdapter, secretKey)

		err := service.Logout(ctx, "invalid-token")
		assert.ErrorContains(t, err, "erro ao validar token")
	})

	t.Run("erro ao adicionar na blacklist", func(t *testing.T) {
		token := generateTestToken(secretKey, time.Now().Add(time.Hour))
		mockBL := new(mockBlacklist)
		mockBL.On("Add", mock.Anything, token, mock.AnythingOfType("time.Duration")).Return(errors.New("falha redis"))

		service := NewLogoutService(mockBL, loggerAdapter, secretKey)
		err := service.Logout(ctx, token)
		assert.ErrorContains(t, err, "erro ao realizar logout")
		mockBL.AssertExpectations(t)
	})

	t.Run("claims inválidas", func(t *testing.T) {
		claims := jwt.MapClaims{
			"foo": "bar",
			// Note que não tem "exp"
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secretKey))

		mockBL := new(mockBlacklist)
		service := NewLogoutService(mockBL, loggerAdapter, secretKey)

		err := service.Logout(ctx, tokenString)
		assert.ErrorContains(t, err, "claim 'exp' ausente ou inválida")
	})

}
