package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateAndValidate_Success(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Minute*5, "auth-service", "store-client")

	tokenStr, err := manager.Generate(123, "user@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := manager.ValidateToken(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, "123", claims["sub"])
	assert.Equal(t, "123", claims["user_id"])
	assert.Equal(t, "user@example.com", claims["email"])
	assert.Equal(t, "auth-service", claims["iss"])
	assert.Equal(t, "store-client", claims["aud"])
}

func TestJWTManager_ValidateToken_Expired(t *testing.T) {
	manager := NewJWTManager("test-secret", -time.Minute*1, "auth-service", "store-client")

	tokenStr, err := manager.Generate(123, "user@example.com")
	assert.NoError(t, err)

	_, err = manager.ValidateToken(tokenStr)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestJWTManager_ValidateToken_InvalidAudience(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Minute*5, "auth-service", "expected-aud")

	// Gera token com audience incorreta
	claims := jwt.MapClaims{
		"sub":     "123",
		"user_id": "123",
		"email":   "user@example.com",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
		"iss":     "auth-service",
		"aud":     "wrong-aud",
		"jti":     "test-id",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("test-secret"))
	assert.NoError(t, err)

	_, err = manager.ValidateToken(tokenStr)
	assert.ErrorIs(t, err, ErrInvalidAudience)
}

func TestJWTManager_ValidateToken_InvalidIssuer(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Minute*5, "expected-iss", "store-client")

	// Gera token com issuer incorreto
	claims := jwt.MapClaims{
		"sub":     "123",
		"user_id": "123",
		"email":   "user@example.com",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
		"iss":     "wrong-issuer",
		"aud":     "store-client",
		"jti":     "test-id",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("test-secret"))
	assert.NoError(t, err)

	_, err = manager.ValidateToken(tokenStr)
	assert.ErrorIs(t, err, ErrInvalidIssuer)
}

func TestJWTManager_ValidateToken_InvalidSignature(t *testing.T) {
	manager := NewJWTManager("correct-secret", time.Minute*5, "auth-service", "store-client")

	// Gera token com outra chave
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     "123",
		"user_id": "123",
		"email":   "user@example.com",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
		"iss":     "auth-service",
		"aud":     "store-client",
		"jti":     "test-id",
	})
	tokenStr, err := token.SignedString([]byte("wrong-secret"))
	assert.NoError(t, err)

	_, err = manager.ValidateToken(tokenStr)
	assert.Error(t, err) // assinatura inválida
}
