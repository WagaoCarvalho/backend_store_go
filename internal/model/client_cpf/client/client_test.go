package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientCpf_Validate(t *testing.T) {

	t.Run("nome em branco", func(t *testing.T) {
		c := &ClientCpf{
			Name:    "",
			Email:   "teste@teste.com",
			CPF:     "12345678909",
			Version: 1,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
		assert.Contains(t, err.Error(), "campo obrigatório")
	})

	t.Run("nome maior que 255 caracteres", func(t *testing.T) {
		longName := strings.Repeat("a", 256)

		c := &ClientCpf{
			Name:    longName,
			Email:   "teste@teste.com",
			CPF:     "12345678909",
			Version: 1,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
		assert.Contains(t, err.Error(), "máximo de 255 caracteres")
	})

	t.Run("email em branco", func(t *testing.T) {
		c := &ClientCpf{
			Name:    "Cliente Teste",
			Email:   "",
			CPF:     "12345678909",
			Version: 1,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
		assert.Contains(t, err.Error(), "campo obrigatório")
	})

	t.Run("email inválido", func(t *testing.T) {
		c := &ClientCpf{
			Name:    "Cliente Teste",
			Email:   "email-invalido",
			CPF:     "12345678909",
			Version: 1,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
		assert.Contains(t, err.Error(), "inválido")
	})

	t.Run("CPF inválido", func(t *testing.T) {
		c := &ClientCpf{
			Name:    "Teste CPF",
			Email:   "teste@teste.com",
			CPF:     "00000000000",
			Version: 1,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF")
		assert.Contains(t, err.Error(), "CPF inválido")
	})

	t.Run("CPF em branco", func(t *testing.T) {
		c := &ClientCpf{
			Name:    "Cliente Teste",
			Email:   "teste@teste.com",
			CPF:     "",
			Version: 1,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF")
		assert.Contains(t, err.Error(), "campo obrigatório")
	})

	t.Run("description maior que 1000 caracteres", func(t *testing.T) {
		longDescription := strings.Repeat("a", 1001)

		c := &ClientCpf{
			Name:        "Cliente Teste",
			Email:       "teste@teste.com",
			CPF:         "12345678909",
			Description: longDescription,
			Version:     1,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Description")
		assert.Contains(t, err.Error(), "máximo de 1000 caracteres")
	})

	t.Run("versão inválida", func(t *testing.T) {
		c := &ClientCpf{
			Name:    "Cliente Teste",
			Email:   "teste@teste.com",
			CPF:     "12345678909",
			Version: 0,
		}

		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version")
	})

	t.Run("cliente válido", func(t *testing.T) {
		c := &ClientCpf{
			Name:    "Cliente Teste",
			Email:   "TESTE@TESTE.COM",
			CPF:     "12345678909",
			Version: 1,
			Status:  true,
		}

		err := c.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "teste@teste.com", c.Email) // normalização obrigatória
	})
}
