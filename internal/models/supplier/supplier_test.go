package models_test

import (
	"strings"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_supplier "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	"github.com/stretchr/testify/assert"
)

func TestSupplier_Validate(t *testing.T) {
	validName := "Fornecedor Exemplo"
	validContact := "Contato Exemplo"
	validCPF := "12345678901"
	validCNPJ := "12345678000199"

	t.Run("Valid supplier with CPF", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CPF:         &validCPF,
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid supplier with CNPJ", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CNPJ:        &validCNPJ,
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("Missing name", func(t *testing.T) {
		s := &models_supplier.Supplier{
			ContactInfo: validContact,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
	})

	t.Run("Missing contact info", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name: validName,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ContactInfo")
	})

	t.Run("Name too long", func(t *testing.T) {
		longName := strings.Repeat("a", 101)
		s := &models_supplier.Supplier{
			Name:        longName,
			ContactInfo: validContact,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
	})

	t.Run("ContactInfo too long", func(t *testing.T) {
		longContact := strings.Repeat("c", 101)
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: longContact,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ContactInfo")
	})

	t.Run("Invalid CPF", func(t *testing.T) {
		invalidCPF := "123"
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CPF:         &invalidCPF,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF")
	})

	t.Run("Invalid CNPJ", func(t *testing.T) {
		invalidCNPJ := "abc"
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CNPJ:        &invalidCNPJ,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CNPJ")
	})

	t.Run("Both CPF and CNPJ filled", func(t *testing.T) {
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CPF:         &validCPF,
			CNPJ:        &validCNPJ,
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF/CNPJ")
	})

	t.Run("Invalid Address", func(t *testing.T) {
		userID := int64(1) // ID obrigatório para passar a validação
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CPF:         &validCPF,
			Address: &models_address.Address{
				UserID: &userID, // Adicionado para passar validação de IDs
				Street: "",      // inválido: campo obrigatório
			},
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Street") // ou outro campo inválido no Address
	})

	t.Run("Invalid Contact", func(t *testing.T) {
		userID := int64(1) // ID obrigatório para passar a validação
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CPF:         &validCPF,
			Contact: &models_contact.Contact{
				UserID:      &userID, // Adicionado para passar validação de IDs
				ContactName: "",      // inválido: campo obrigatório
			},
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ContactName") // ou outro campo inválido no Contact
	})

	t.Run("Valid Address and Contact", func(t *testing.T) {
		userID := int64(1) // ID obrigatório para passar a validação
		s := &models_supplier.Supplier{
			Name:        validName,
			ContactInfo: validContact,
			CNPJ:        &validCNPJ,
			Address: &models_address.Address{
				UserID:     &userID, // Adicionado para passar validação de IDs
				Street:     "Rua X",
				City:       "Cidade",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345-678",
			},
			Contact: &models_contact.Contact{
				UserID:      &userID, // Adicionado para passar validação de IDs
				ContactName: "João",
			},
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

}
