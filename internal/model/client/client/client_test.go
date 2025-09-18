package model

import (
	"strings"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestClient_Validate(t *testing.T) {
	t.Run("nome em branco", func(t *testing.T) {
		c := &Client{
			Name: "",

			CPF: utils.StrToPtr("12345678909"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
		assert.Contains(t, err.Error(), "campo obrigatório")
	})

	t.Run("nome maior que 255 caracteres", func(t *testing.T) {
		longName := strings.Repeat("a", 256)
		c := &Client{
			Name: longName,

			CPF: utils.StrToPtr("12345678909"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
		assert.Contains(t, err.Error(), "máximo de 255 caracteres")
	})

	t.Run("CPF e CNPJ preenchidos simultaneamente", func(t *testing.T) {
		c := &Client{
			Name: "Teste",

			CPF:  utils.StrToPtr("12345678909"),
			CNPJ: utils.StrToPtr("45384524888010"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF/CNPJ")
	})

	t.Run("CPF inválido", func(t *testing.T) {
		c := &Client{
			Name: "Teste CPF",

			CPF: utils.StrToPtr("00000000000"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF inválido")
	})

	t.Run("CNPJ inválido", func(t *testing.T) {
		c := &Client{
			Name: "Teste CNPJ",
			CNPJ: utils.StrToPtr("00000000000000"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CNPJ inválido")
	})

	t.Run("sem CPF e sem CNPJ", func(t *testing.T) {
		c := &Client{
			Name: "Cliente Teste",
			// CPF e CNPJ nulos
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF/CNPJ")
		assert.Contains(t, err.Error(), "deve informar CPF ou CNPJ")
	})

	t.Run("sucesso PF válido", func(t *testing.T) {
		c := &Client{
			Name: "Cliente PF",

			CPF:  utils.StrToPtr("12345678909"),
			CNPJ: nil,
		}
		err := c.Validate()
		assert.NoError(t, err)
	})

	t.Run("sucesso PJ válido", func(t *testing.T) {
		c := &Client{
			Name: "Cliente PJ",
			CNPJ: utils.StrToPtr("45384524888010"),
			CPF:  nil,
		}
		err := c.Validate()
		assert.NoError(t, err)
	})
}
