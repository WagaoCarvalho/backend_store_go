package models_test

import (
	"strings"
	"testing"

	models_supplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier"
	"github.com/stretchr/testify/assert"
)

func TestSupplier_Validate(t *testing.T) {
	validName := "Fornecedor Exemplo"
	validCPF := "12345678901"
	validCNPJ := "12345678000199"

	t.Run("Name em branco", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name: "",
			CPF:  &validCPF,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
		assert.Contains(t, err.Error(), "obrigatório")
	})

	t.Run("Valid supplier with CPF", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name: validName,
			CPF:  &validCPF,
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid supplier with CNPJ", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name: validName,
			CNPJ: &validCNPJ,
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("Name too long", func(t *testing.T) {
		longName := strings.Repeat("a", 101)
		s := &models_supplier.Supplier{
			Name: longName,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
	})

	t.Run("Invalid CPF", func(t *testing.T) {
		invalidCPF := "123"
		s := &models_supplier.Supplier{
			Name: validName,
			CPF:  &invalidCPF,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF")
	})

	t.Run("Invalid CNPJ", func(t *testing.T) {
		invalidCNPJ := "abc"
		s := &models_supplier.Supplier{
			Name: validName,
			CNPJ: &invalidCNPJ,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CNPJ")
	})

	t.Run("Both CPF and CNPJ filled", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name: validName,
			CPF:  &validCPF,
			CNPJ: &validCNPJ,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF/CNPJ")
	})

}
