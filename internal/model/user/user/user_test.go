package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	validEmail := "user@example.com"
	validUsername := "usuarioValido"
	validPassword := "Senha123"

	t.Run("Missing Username", func(t *testing.T) {
		u := &User{
			Username: "",
			Email:    validEmail,
			Password: validPassword,
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Username")
	})

	t.Run("Username too short", func(t *testing.T) {
		u := &User{
			Username: "ab",
			Email:    validEmail,
			Password: validPassword,
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "3 e 50")
	})

	t.Run("Username too long", func(t *testing.T) {
		u := &User{
			Username: "a" + string(make([]byte, 51)), // 52 chars
			Email:    validEmail,
			Password: validPassword,
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "3 e 50")
	})

	t.Run("Missing Email", func(t *testing.T) {
		u := &User{
			Username: validUsername,
			Email:    "",
			Password: validPassword,
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		u := &User{
			Username: validUsername,
			Email:    "invalid-email",
			Password: validPassword,
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
	})

	t.Run("Email too long", func(t *testing.T) {
		longEmail := "a" + string(make([]byte, 100)) + "@example.com" // > 100 chars
		u := &User{
			Username: validUsername,
			Email:    longEmail,
			Password: validPassword,
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "máximo de 100")
	})

	t.Run("Missing Password", func(t *testing.T) {
		u := &User{
			Username: validUsername,
			Email:    validEmail,
			Password: "",
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Password")
	})

	t.Run("Password too short", func(t *testing.T) {
		u := &User{
			Username: validUsername,
			Email:    validEmail,
			Password: "aB1", // menos de 8 chars
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mínimo de 8")
	})

	t.Run("Password missing complexity", func(t *testing.T) {
		// Test password without uppercase letter
		u := &User{
			Username: validUsername,
			Email:    validEmail,
			Password: "senha1234",
		}

		err := u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deve conter")

		// Test password without lowercase letter
		u.Password = "SENHA1234"
		err = u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deve conter")

		// Test password without number
		u.Password = "SenhaSenha"
		err = u.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deve conter")
	})

	t.Run("should return nil when all fields are valid", func(t *testing.T) {
		user := &User{
			Username: "validuser",
			Email:    "valid@example.com",
			Password: "Strong123", // Assume que hasPasswordComplexity valida corretamente esse padrão
		}

		err := user.Validate()

		assert.NoError(t, err)
	})

}
