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
			Name:       "",
			ClientType: "PF",
			CPF:        utils.StrToPtr("12345678909"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
		assert.Contains(t, err.Error(), "campo obrigatório")
	})

	t.Run("nome maior que 255 caracteres", func(t *testing.T) {
		longName := strings.Repeat("a", 256)
		c := &Client{
			Name:       longName,
			ClientType: "PF",
			CPF:        utils.StrToPtr("12345678909"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
		assert.Contains(t, err.Error(), "máximo de 255 caracteres")
	})

	t.Run("ClientType em branco", func(t *testing.T) {
		c := &Client{
			Name: "Cliente Teste",
			CPF:  utils.StrToPtr("12345678909"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ClientType")
		assert.Contains(t, err.Error(), "campo obrigatório")
	})

	t.Run("ClientType inválido", func(t *testing.T) {
		c := &Client{
			Name:       "Teste",
			ClientType: "Outro",
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ClientType")
	})

	t.Run("CPF obrigatório para PF", func(t *testing.T) {
		c := &Client{
			Name:       "Cliente PF",
			ClientType: "PF",
			CPF:        nil,
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF")
		assert.Contains(t, err.Error(), "obrigatório para pessoa física")
	})

	t.Run("CNPJ obrigatório para PJ", func(t *testing.T) {
		c := &Client{
			Name:       "Cliente PJ",
			ClientType: "PJ",
			CNPJ:       nil,
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CNPJ")
		assert.Contains(t, err.Error(), "obrigatório para pessoa jurídica")
	})

	t.Run("CPF e CNPJ preenchidos simultaneamente", func(t *testing.T) {
		c := &Client{
			Name:       "Teste",
			ClientType: "PF",
			CPF:        utils.StrToPtr("12345678909"),
			CNPJ:       utils.StrToPtr("45384524888010"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF/CNPJ")
	})

	t.Run("CPF inválido", func(t *testing.T) {
		c := &Client{
			Name:       "Teste CPF",
			ClientType: "PF",
			CPF:        utils.StrToPtr("00000000000"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF inválido")
	})

	t.Run("CNPJ inválido", func(t *testing.T) {
		c := &Client{
			Name:       "Teste CNPJ",
			ClientType: "PJ",
			CNPJ:       utils.StrToPtr("00000000000000"),
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CNPJ inválido")
	})

	t.Run("sucesso PF válido", func(t *testing.T) {
		c := &Client{
			Name:       "Cliente PF",
			ClientType: "PF",
			CPF:        utils.StrToPtr("12345678909"),
			CNPJ:       nil,
		}
		err := c.Validate()
		assert.NoError(t, err)
	})

	t.Run("sucesso PJ válido", func(t *testing.T) {
		c := &Client{
			Name:       "Cliente PJ",
			ClientType: "PJ",
			CNPJ:       utils.StrToPtr("45384524888010"),
			CPF:        nil,
		}
		err := c.Validate()
		assert.NoError(t, err)
	})
}
