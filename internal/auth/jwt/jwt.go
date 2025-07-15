package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
		"user_id": strconv.FormatInt(uid, 10), // <- obrigatório como string
		"email":   email,
		"exp":     time.Now().Add(j.TokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}
