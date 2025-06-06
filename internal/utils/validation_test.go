package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilsFunctions(t *testing.T) {
	t.Run("IsStringEmpty", func(t *testing.T) {
		assert.True(t, IsStringEmpty(""), "string vazia deve retornar true")
		assert.False(t, IsStringEmpty(" "), "string com espaço não deve retornar true")
		assert.False(t, IsStringEmpty("abc"), "string não vazia deve retornar false")
	})

	t.Run("IsValidEmail", func(t *testing.T) {
		validEmails := []string{
			"user@example.com",
			"user.name+tag+sorting@example.co.uk",
			"user_name@example.io",
			"user-name@example.com",
			"1234567890@example.com",
		}

		invalidEmails := []string{
			"",
			"userexample.com",
			"user@.com",
			"user@com",
			"user@site..com",
			"user@site,com",
			"user@site@site.com",
			"@example.com",
		}

		for _, email := range validEmails {
			assert.True(t, IsValidEmail(email), "email válido deve retornar true: %s", email)
		}

		for _, email := range invalidEmails {
			assert.False(t, IsValidEmail(email), "email inválido deve retornar false: %s", email)
		}
	})
}
