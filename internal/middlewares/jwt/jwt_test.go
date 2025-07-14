package middlewares

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
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

// Copiar a mesma chave do middleware para testes
type contextKey_test string

const userClaimsKey_test = contextKey("user")

func generateValidToken(t *testing.T, secret string, expiration time.Duration) string {
	t.Helper()
	claims := jwt.MapClaims{
		"uid":   1,
		"email": "test@example.com",
		"exp":   time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("erro ao gerar token: %v", err)
	}
	return signed
}

func TestIsAuthByBearerToken(t *testing.T) {
	log := logrus.New()
	log.SetOutput(io.Discard)
	loggerAdapter := logger.NewLoggerAdapter(log)
	secret := "test-secret"

	t.Run("token ausente retorna 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token ausente")
	})

	t.Run("formato de token inválido retorna 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "TokenSemBearer")
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Formato de token inválido")
	})

	t.Run("token revogado retorna 401", func(t *testing.T) {
		token := generateValidToken(t, secret, time.Minute)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, token).Return(true, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token revogado")
		blacklist.AssertExpectations(t)
	})

	t.Run("token inválido retorna 401", func(t *testing.T) {
		invalidToken := "invalid.token.value"
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+invalidToken)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, invalidToken).Return(false, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token inválido")
		blacklist.AssertExpectations(t)
	})

	t.Run("token válido preenche contexto e passa para next handler", func(t *testing.T) {
		validToken := generateValidToken(t, secret, time.Minute)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, validToken).Return(false, nil)

		// Handler que lê o contexto com a mesma chave usada no middleware
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(userClaimsKey)
			if claims == nil {
				http.Error(w, "claims ausentes", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		blacklist.AssertExpectations(t)
	})

	t.Run("token com método de assinatura inválido retorna 401", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("erro ao gerar chave RSA: %v", err)
		}

		claims := jwt.MapClaims{
			"uid":   1,
			"email": "test@example.com",
			"exp":   time.Now().Add(time.Minute).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		tokenString, err := token.SignedString(privateKey)
		if err != nil {
			t.Fatalf("erro ao assinar token RS256: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, tokenString).Return(false, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token inválido")
		blacklist.AssertExpectations(t)
	})

	t.Run("erro ao consultar blacklist retorna 500", func(t *testing.T) {
		token := generateValidToken(t, secret, time.Minute)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, token).Return(false, assert.AnError)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Erro interno de autenticação")
		blacklist.AssertExpectations(t)
	})

	t.Run("token expirado retorna 401", func(t *testing.T) {
		expiredToken := generateValidToken(t, secret, -time.Minute)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, expiredToken).Return(false, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token expirado")
		blacklist.AssertExpectations(t)
	})

	t.Run("assinatura inválida retorna 401", func(t *testing.T) {
		otherSecret := "wrong-secret"
		invalidSigToken := generateValidToken(t, otherSecret, time.Minute)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+invalidSigToken)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, invalidSigToken).Return(false, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Assinatura inválida")
		blacklist.AssertExpectations(t)
	})

	t.Run("token expirado com parse válido retorna 401", func(t *testing.T) {
		expiredClaims := jwt.MapClaims{
			"uid":   1,
			"email": "test@example.com",
			"exp":   time.Now().Add(-time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("erro ao assinar token expirado: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, tokenString).Return(false, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token expirado")
		blacklist.AssertExpectations(t)
	})

	t.Run("token com exp ausente ou inválido retorna 401", func(t *testing.T) {
		claims := jwt.MapClaims{
			"uid":   1,
			"email": "test@example.com",
			// exp ausente intencionalmente
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("erro ao gerar token sem exp: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()

		blacklist := new(mockBlacklist)
		blacklist.On("IsBlacklisted", mock.Anything, tokenString).Return(false, nil)

		// Handler que não deve ser chamado
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := IsAuthByBearerToken(blacklist, loggerAdapter, secret)
		middleware(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Token inválido")
		blacklist.AssertExpectations(t)
	})
}
