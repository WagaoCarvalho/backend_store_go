package utils

import (
	"fmt"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(uid int64, email string) (string, error) {
	secretKey := []byte(config.LoadJwtConfig().SecretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   uid,
		"email": email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar token: %w", err)
	}

	return tokenString, nil
}
