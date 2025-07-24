package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
)

func TestUserFull_Validate(t *testing.T) {
	t.Run("deve retornar erro se User for nil", func(t *testing.T) {
		uf := &UserFull{
			User:       nil,
			Address:    &models_address.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    &models_contact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []models_user_categories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "usuário é obrigatório")
	})

	t.Run("deve retornar erro se Address for nil", func(t *testing.T) {
		uf := &UserFull{
			User:       &models_user.User{Username: "usuario"},
			Address:    nil,
			Contact:    &models_contact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []models_user_categories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "endereço é obrigatório")
	})

	t.Run("deve retornar erro se Contact for nil", func(t *testing.T) {
		uf := &UserFull{
			User:       &models_user.User{Username: "usuario"},
			Address:    &models_address.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    nil,
			Categories: []models_user_categories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "contato é obrigatório")
	})

	t.Run("deve retornar erro se Categories estiver vazio", func(t *testing.T) {
		uf := &UserFull{
			User:       &models_user.User{Username: "usuario"},
			Address:    &models_address.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    &models_contact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []models_user_categories.UserCategory{},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "pelo menos uma categoria é obrigatória")
	})

	t.Run("deve passar quando todos os campos obrigatórios estiverem presentes", func(t *testing.T) {
		uf := &UserFull{
			User:       &models_user.User{Username: "usuario", Email: "usuario@example.com", Password: "Senha123"},
			Address:    &models_address.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    &models_contact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []models_user_categories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.NoError(t, err)
	})
}
