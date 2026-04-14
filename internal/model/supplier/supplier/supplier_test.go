package model

import (
	"strings"
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestSupplier_Validate(t *testing.T) {
	validName := "Fornecedor Exemplo"
	validCPF := "12345678901"
	validCNPJ := "12345678000199"

	t.Run("Name em branco", func(t *testing.T) {
		s := &Supplier{
			Name: "",
			CPF:  &validCPF,
		}
		err := s.Validate()
		assert.Error(t, err)

		validationErrs, ok := err.(validators.ValidationErrors)
		assert.True(t, ok)
		assert.Len(t, validationErrs, 1)
		assert.Equal(t, "name", validationErrs[0].Field)
		assert.Equal(t, validators.MsgRequiredField, validationErrs[0].Message)
	})

	t.Run("Valid supplier with CPF", func(t *testing.T) {
		s := &Supplier{
			Name: validName,
			CPF:  &validCPF,
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid supplier with CNPJ", func(t *testing.T) {
		s := &Supplier{
			Name: validName,
			CNPJ: &validCNPJ,
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("Name too long", func(t *testing.T) {
		longName := strings.Repeat("a", 101)
		s := &Supplier{
			Name: longName,
		}
		err := s.Validate()
		assert.Error(t, err)

		validationErrs, ok := err.(validators.ValidationErrors)
		assert.True(t, ok)
		assert.Len(t, validationErrs, 1)
		assert.Equal(t, "name", validationErrs[0].Field)
		assert.Equal(t, validators.MsgMax100, validationErrs[0].Message)
	})

	t.Run("Invalid CPF", func(t *testing.T) {
		invalidCPF := "123"
		s := &Supplier{
			Name: validName,
			CPF:  &invalidCPF,
		}
		err := s.Validate()
		assert.Error(t, err)

		validationErrs, ok := err.(validators.ValidationErrors)
		assert.True(t, ok)
		assert.Len(t, validationErrs, 1)
		assert.Equal(t, "cpf", validationErrs[0].Field)
	})

	t.Run("Invalid CNPJ", func(t *testing.T) {
		invalidCNPJ := "abc"
		s := &Supplier{
			Name: validName,
			CNPJ: &invalidCNPJ,
		}
		err := s.Validate()
		assert.Error(t, err)

		validationErrs, ok := err.(validators.ValidationErrors)
		assert.True(t, ok)
		assert.Len(t, validationErrs, 1)
		assert.Equal(t, "cnpj", validationErrs[0].Field)
	})

	t.Run("Both CPF and CNPJ filled", func(t *testing.T) {
		s := &Supplier{
			Name: validName,
			CPF:  &validCPF,
			CNPJ: &validCNPJ,
		}
		err := s.Validate()
		assert.Error(t, err)

		validationErrs, ok := err.(validators.ValidationErrors)
		assert.True(t, ok)
		assert.Len(t, validationErrs, 1)
		assert.Equal(t, "cpf_cnpj", validationErrs[0].Field)
		assert.Equal(t, validators.MsgInvalidAssociation, validationErrs[0].Message)
	})

	t.Run("Multiple validation errors", func(t *testing.T) {
		s := &Supplier{
			Name: "", // Name vazio
			CPF:  &validCPF,
			CNPJ: &validCNPJ, // Ambos preenchidos
		}
		err := s.Validate()
		assert.Error(t, err)

		validationErrs, ok := err.(validators.ValidationErrors)
		assert.True(t, ok)
		// Deve ter pelo menos 2 erros: Name vazio + ambos preenchidos
		assert.GreaterOrEqual(t, len(validationErrs), 2)

		// Verifica se contém o erro de name
		hasNameError := false
		hasBothError := false
		for _, ve := range validationErrs {
			if ve.Field == "name" {
				hasNameError = true
				assert.Equal(t, validators.MsgRequiredField, ve.Message)
			}
			if ve.Field == "cpf_cnpj" {
				hasBothError = true
				assert.Equal(t, validators.MsgInvalidAssociation, ve.Message)
			}
		}
		assert.True(t, hasNameError, "Deveria ter erro no campo name")
		assert.True(t, hasBothError, "Deveria ter erro no campo cpf_cnpj")
	})
}

// Teste adicional para verificar a mensagem de erro formatada corretamente
func TestSupplier_Validate_ErrorMessage(t *testing.T) {
	s := &Supplier{
		Name: "",
		CPF:  nil,
		CNPJ: nil,
	}
	err := s.Validate()
	assert.Error(t, err)

	// Verifica se o erro implementa a interface Error
	errorString := err.Error()
	assert.NotEmpty(t, errorString)

	// Deve conter as informações dos erros
	assert.Contains(t, errorString, "name")
}
