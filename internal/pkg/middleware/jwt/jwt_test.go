package middleware

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBlacklist struct {
	mock.Mock
}

func (m *mockBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

type mockJWTService struct {
	mock.Mock
}

func (m *mockJWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

func (m *mockJWTService) Generate(uid int64, email string) (string, error) {
	args := m.Called(uid, email)
	return args.String(0), args.Error(1)
}

func buildJWT(t *testing.T, manager *auth.JWTManager, duration time.Duration) string {
	manager.TokenDuration = duration
	token, err := manager.Generate(1, "test@example.com")
	assert.NoError(t, err)
	return token
}

func TestIsAuthByBearerToken(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	loggerAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("sem header Authorization", func(t *testing.T) {
		mockJWT := new(mockJWTService)
		mockBL := new(mockBlacklist)

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		middleware := IsAuthByBearerToken(mockBL, loggerAdapter, mockJWT)
		middleware(handler).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, rec.Body.String(), auth.ErrTokenMissing.Error())
	})

	t.Run("formato inválido do header", func(t *testing.T) {
		mockJWT := new(mockJWTService)
		mockBL := new(mockBlacklist)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidFormatToken")
		rec := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		middleware := IsAuthByBearerToken(mockBL, loggerAdapter, mockJWT)
		middleware(handler).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, rec.Body.String(), auth.ErrTokenInvalidFormat.Error())
	})

	t.Run("token revogado", func(t *testing.T) {
		mockJWT := new(mockJWTService)
		mockBL := new(mockBlacklist)

		token := "fake-token"
		mockBL.On("IsBlacklisted", mock.Anything, token).Return(true, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		middleware := IsAuthByBearerToken(mockBL, loggerAdapter, mockJWT)
		middleware(handler).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, rec.Body.String(), auth.ErrTokenRevoked.Error())
	})

	t.Run("erro interno na blacklist", func(t *testing.T) {
		mockJWT := new(mockJWTService)
		mockBL := new(mockBlacklist)

		token := "fake-token"
		mockBL.On("IsBlacklisted", mock.Anything, token).Return(false, errors.New("internal"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		middleware := IsAuthByBearerToken(mockBL, loggerAdapter, mockJWT)
		middleware(handler).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Contains(t, rec.Body.String(), auth.ErrInternalAuth.Error())
	})

	t.Run("token expirado", func(t *testing.T) {
		duration := -1 * time.Minute
		jwtManager := auth.NewJWTManager("test-key", duration, "auth-service", "store-client")
		token := buildJWT(t, jwtManager, duration)

		mockJWT := new(mockJWTService)
		mockBL := new(mockBlacklist)
		mockBL.On("IsBlacklisted", mock.Anything, token).Return(false, nil)

		// Forçar retorno de erro de expiração com ValidationError correto
		mockJWT.On("ValidateToken", token).
			Return(jwt.MapClaims{}, auth.ErrTokenExpired)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		middleware := IsAuthByBearerToken(mockBL, loggerAdapter, mockJWT)
		middleware(handler).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, rec.Body.String(), auth.ErrTokenExpired.Error())
	})

	t.Run("claim user_id ausente ou inválida", func(t *testing.T) {
		duration := 5 * time.Minute
		jwtManager := auth.NewJWTManager("test-key", duration, "auth-service", "store-client")
		token := buildJWT(t, jwtManager, duration)

		mockJWT := new(mockJWTService)
		mockBL := new(mockBlacklist)

		mockBL.On("IsBlacklisted", mock.Anything, token).Return(false, nil)

		// Retorna claims sem "user_id" ou com valor inválido
		mockJWT.On("ValidateToken", token).Return(jwt.MapClaims{
			"email": "test@example.com",
			"exp":   float64(time.Now().Add(duration).Unix()),
			// "user_id" omitido de propósito
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		loggerAdapter := logger.NewLoggerAdapter(logrus.New())

		middleware := IsAuthByBearerToken(mockBL, loggerAdapter, mockJWT)
		middleware(handler).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, rec.Body.String(), auth.ErrTokenInvalid.Error())
	})

	t.Run("token válido", func(t *testing.T) {
		duration := 5 * time.Minute
		jwtManager := auth.NewJWTManager("test-key", duration, "auth-service", "store-client")
		token := buildJWT(t, jwtManager, duration)

		mockJWT := new(mockJWTService)
		mockBL := new(mockBlacklist)

		mockBL.On("IsBlacklisted", mock.Anything, token).Return(false, nil)

		mockJWT.On("ValidateToken", token).Return(jwt.MapClaims{
			"user_id": "1", // string, conforme esperado pelo middleware
			"email":   "test@example.com",
			"exp":     float64(time.Now().Add(duration).Unix()),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		loggerAdapter := logger.NewLoggerAdapter(logrus.New())

		middleware := IsAuthByBearerToken(mockBL, loggerAdapter, mockJWT)
		middleware(handler).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, rec.Body.String(), "ok")
	})

}
