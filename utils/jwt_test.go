package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	uid := int64(123)
	email := "user@example.com"

	t.Run("Gera token sem erro", func(t *testing.T) {
		tokenString, err := GenerateJWT(uid, email)
		assert.NoError(t, err, "não deve retornar erro ao gerar token")
		assert.NotEmpty(t, tokenString, "token gerado não deve ser vazio")
		assert.Equal(t, 2, strings.Count(tokenString, "."), "token deve ter 3 partes separadas por '.'")
	})

	t.Run("Token contém claims corretos", func(t *testing.T) {
		tokenString, err := GenerateJWT(uid, email)
		assert.NoError(t, err)

		secretKey := []byte(config.LoadJwtConfig().SecretKey)

		// Parseia token para verificar claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid, "token deve ser válido")

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok, "claims devem ser do tipo MapClaims")

		assert.Equal(t, float64(uid), claims["uid"], "uid deve ser igual ao passado")
		assert.Equal(t, email, claims["email"], "email deve ser igual ao passado")

		exp, ok := claims["exp"].(float64)
		assert.True(t, ok, "exp deve existir e ser float64")

		assert.Greater(t, int64(exp), time.Now().Unix(), "exp deve ser timestamp futuro")
	})
}
