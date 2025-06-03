package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ✅ Esta struct pertence ao jwt.go
type JWTManager struct {
	SecretKey     string
	TokenDuration time.Duration
}

type JWTGenerator interface {
	Generate(uid int64, email string) (string, error)
}

func NewJWTManager(secretKey string, duration time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey:     secretKey,
		TokenDuration: duration,
	}
}

func (j *JWTManager) Generate(uid int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"uid":   uid,
		"email": email,
		"exp":   time.Now().Add(j.TokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}
