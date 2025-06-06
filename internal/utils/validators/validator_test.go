package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilsFunctions(t *testing.T) {
	t.Run("IsBlank", func(t *testing.T) {
		assert.True(t, IsBlank(""), "string vazia deve retornar true")
		assert.True(t, IsBlank("   "), "string com apenas espaços deve retornar true")
		assert.False(t, IsBlank("abc"), "string não vazia deve retornar false")
		assert.False(t, IsBlank(" abc "), "string com conteúdo e espaços deve retornar false")
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

	t.Run("IsValidPostalCode", func(t *testing.T) {
		validCodes := []string{"12345-678", "12345678"}
		invalidCodes := []string{"1234-567", "1234567", "abcde-fgh"}

		for _, code := range validCodes {
			assert.True(t, IsValidPostalCode(code), "CEP válido deve retornar true: %s", code)
		}
		for _, code := range invalidCodes {
			assert.False(t, IsValidPostalCode(code), "CEP inválido deve retornar false: %s", code)
		}
	})

	t.Run("IsValidBrazilianState", func(t *testing.T) {
		validStates := []string{"SP", "RJ", "MG", "DF", "AC"}
		invalidStates := []string{"XX", "ABC", "", "sp", "rj"}

		for _, state := range validStates {
			assert.True(t, IsValidBrazilianState(state), "Estado válido deve retornar true: %s", state)
		}
		for _, state := range invalidStates {
			assert.False(t, IsValidBrazilianState(state), "Estado inválido deve retornar false: %s", state)
		}
	})

	t.Run("IsValidPhone", func(t *testing.T) {
		validPhones := []string{"(11) 1234-5678"}
		invalidPhones := []string{"11 1234-5678", "(11)1234-5678", "(11) 12345-6789", "1234-5678"}

		for _, phone := range validPhones {
			assert.True(t, IsValidPhone(phone), "Telefone válido deve retornar true: %s", phone)
		}
		for _, phone := range invalidPhones {
			assert.False(t, IsValidPhone(phone), "Telefone inválido deve retornar false: %s", phone)
		}
	})

	t.Run("IsValidCell", func(t *testing.T) {
		validCells := []string{"(11) 91234-5678"}
		invalidCells := []string{"(11) 1234-5678", "11 91234-5678", "(11)91234-5678", "(11) 912345678"}

		for _, cell := range validCells {
			assert.True(t, IsValidCell(cell), "Celular válido deve retornar true: %s", cell)
		}
		for _, cell := range invalidCells {
			assert.False(t, IsValidCell(cell), "Celular inválido deve retornar false: %s", cell)
		}
	})
}
