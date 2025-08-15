package utils

import (
	"strings"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/pkg/utils"
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
		validCodes := []string{"12345678"}
		invalidCodes := []string{"1234-567", "1234567", "abcde-fgh"}

		for _, code := range validCodes {
			assert.True(t, IsValidPostalCode(code), "CEP válido deve retornar true: %s", code)
		}
		for _, code := range invalidCodes {
			assert.False(t, IsValidPostalCode(code), "CEP inválido deve retornar false: %s", code)
		}
	})

	t.Run("IsValidPostalCode retorna false para códigos inválidos", func(t *testing.T) {
		invalidPostalCodes := []string{
			"00000-000", // explícito no if
			"00000000",  // explícito no if
			"1234-567",  // formato inválido
			"1234567",   // formato inválido
			"abcde-fgh", // caracteres inválidos
			"",          // vazio
		}

		for _, code := range invalidPostalCodes {
			assert.False(t, IsValidPostalCode(code), "Código postal inválido deve retornar false: %s", code)
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
		validPhones := []string{
			"(11) 1234-5678",
			"(11)1234-5678",
			"11 1234 5678",
			"1112345678",
		}
		invalidPhones := []string{
			"(11) 12345-6789", // celular (11 dígitos)
			"1234-5678",       // sem DDD (8 dígitos)
			"111234567",       // 9 dígitos
			"abcdefghij",      // não numérico
		}

		for _, phone := range validPhones {
			assert.True(t, IsValidPhone(phone), "Telefone válido deve retornar true: %s", phone)
		}
		for _, phone := range invalidPhones {
			assert.False(t, IsValidPhone(phone), "Telefone inválido deve retornar false: %s", phone)
		}
	})

	t.Run("IsValidCell", func(t *testing.T) {
		validCells := []string{
			"(11) 91234-5678",
			"(11)91234-5678",
			"11 91234 5678",
			"11912345678",
		}
		invalidCells := []string{
			"(11) 1234-5678", // fixo (10 dígitos)
			"1191234567",     // 10 dígitos
			"912345678",      // 9 dígitos
			"abcdefghijk",    // não numérico
		}

		for _, cell := range validCells {
			assert.True(t, IsValidCell(cell), "Celular válido deve retornar true: %s", cell)
		}
		for _, cell := range invalidCells {
			assert.False(t, IsValidCell(cell), "Celular inválido deve retornar false: %s", cell)
		}
	})

	t.Run("IsValidCPF", func(t *testing.T) {
		validCPFs := []string{
			"12345678901",
			"98765432109",
		}
		invalidCPFs := []string{
			"00000000000",
			"1234567890",     // 10 dígitos
			"123456789012",   // 12 dígitos
			"123.456.789-01", // com pontuação
			"abcdefghi01",    // letras
			"",               // vazio
		}

		for _, cpf := range validCPFs {
			assert.True(t, IsValidCPF(cpf), "CPF válido deve retornar true: %s", cpf)
		}
		for _, cpf := range invalidCPFs {
			assert.False(t, IsValidCPF(cpf), "CPF inválido deve retornar false: %s", cpf)
		}
	})

	t.Run("IsValidCNPJ", func(t *testing.T) {
		validCNPJs := []string{
			"12345678000199",
			"11222333000181",
		}
		invalidCNPJs := []string{
			"00000000000000",
			"1234567800019",      // 13 dígitos
			"123456780001999",    // 15 dígitos
			"12.345.678/0001-99", // com pontuação
			"abcd5678000199",     // letras
			"",                   // vazio
		}

		for _, cnpj := range validCNPJs {
			assert.True(t, IsValidCNPJ(cnpj), "CNPJ válido deve retornar true: %s", cnpj)
		}
		for _, cnpj := range invalidCNPJs {
			assert.False(t, IsValidCNPJ(cnpj), "CNPJ inválido deve retornar false: %s", cnpj)
		}
	})

}

func TestValidateSingleNonNil(t *testing.T) {
	var (
		a int64 = 1
		b int64 = 2
	)

	tests := []struct {
		name     string
		input    []*int64
		expected bool
	}{
		{
			name:     "nenhum valor",
			input:    []*int64{},
			expected: false,
		},
		{
			name:     "todos nil",
			input:    []*int64{nil, nil},
			expected: false,
		},
		{
			name:     "um valor não-nil",
			input:    []*int64{&a, nil, nil},
			expected: true,
		},
		{
			name:     "dois valores não-nil",
			input:    []*int64{&a, &b},
			expected: false,
		},
		{
			name:     "três valores não-nil",
			input:    []*int64{&a, &b, &a},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSingleNonNil(tt.input...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateStrongPassword(t *testing.T) {
	t.Run("campo obrigatório", func(t *testing.T) {
		err := ValidateStrongPassword("")
		assert.Error(t, err)
		ve, ok := err.(*utils.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "Password", ve.Field)
		assert.Equal(t, "campo obrigatório", ve.Message)
	})

	t.Run("senha muito curta", func(t *testing.T) {
		err := ValidateStrongPassword("Ab1@")
		assert.Error(t, err)
		ve, ok := err.(*utils.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "mínimo de 8 caracteres", ve.Message)
	})

	t.Run("senha muito longa", func(t *testing.T) {
		longPwd := strings.Repeat("A1a@", 20) // 80 caracteres
		err := ValidateStrongPassword(longPwd)
		assert.Error(t, err)
		ve, ok := err.(*utils.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "máximo de 64 caracteres", ve.Message)
	})

	t.Run("senha sem maiúscula", func(t *testing.T) {
		err := ValidateStrongPassword("abcdef1@")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "senha deve conter letras maiúsculas")
	})

	t.Run("senha sem minúscula", func(t *testing.T) {
		err := ValidateStrongPassword("ABCDEF1@")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "senha deve conter letras maiúsculas, minúsculas, números e caracteres especiais")
	})

	t.Run("senha sem número", func(t *testing.T) {
		err := ValidateStrongPassword("Abcdefg@")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "senha deve conter letras maiúsculas, minúsculas, números e caracteres especiais")
	})

	t.Run("senha sem símbolo", func(t *testing.T) {
		err := ValidateStrongPassword("Abcdef12")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "senha deve conter letras maiúsculas, minúsculas, números e caracteres especiais")
	})

	t.Run("senha válida", func(t *testing.T) {
		err := ValidateStrongPassword("Abcdef1@")
		assert.NoError(t, err)
	})
}
