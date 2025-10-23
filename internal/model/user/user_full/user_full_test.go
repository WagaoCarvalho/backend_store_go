package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	modelsAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelsContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	modelsUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category"
)

func TestUserFull_Validate(t *testing.T) {
	t.Run("deve retornar erro se User for nil", func(t *testing.T) {
		uf := &UserFull{
			User:       nil,
			Address:    &modelsAddress.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    &modelsContact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []modelsUserCategories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "usuário é obrigatório")
	})

	t.Run("deve retornar erro se Address for nil", func(t *testing.T) {
		uf := &UserFull{
			User:       &modelsUser.User{Username: "usuario"},
			Address:    nil,
			Contact:    &modelsContact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []modelsUserCategories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "endereço é obrigatório")
	})

	t.Run("deve retornar erro se Contact for nil", func(t *testing.T) {
		uf := &UserFull{
			User:       &modelsUser.User{Username: "usuario"},
			Address:    &modelsAddress.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    nil,
			Categories: []modelsUserCategories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "contato é obrigatório")
	})

	t.Run("deve retornar erro se Categories estiver vazio", func(t *testing.T) {
		uf := &UserFull{
			User:       &modelsUser.User{Username: "usuario"},
			Address:    &modelsAddress.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    &modelsContact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []modelsUserCategories.UserCategory{},
		}
		err := uf.Validate()
		assert.EqualError(t, err, "pelo menos uma categoria é obrigatória")
	})

	t.Run("deve passar quando todos os campos obrigatórios estiverem presentes", func(t *testing.T) {
		uf := &UserFull{
			User:       &modelsUser.User{Username: "usuario", Email: "usuario@example.com", Password: "Senha123"},
			Address:    &modelsAddress.Address{Street: "Rua A", City: "Cidade B", State: "SP", Country: "Brasil", PostalCode: "99999999"},
			Contact:    &modelsContact.Contact{ContactName: "João da Silva", Phone: "1199999999", Email: "contato@example.com"},
			Categories: []modelsUserCategories.UserCategory{{ID: 1}},
		}
		err := uf.Validate()
		assert.NoError(t, err)
	})
}
