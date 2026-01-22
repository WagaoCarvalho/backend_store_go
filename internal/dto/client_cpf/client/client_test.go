package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
)

func TestToClientCpfModel(t *testing.T) {
	dtoInput := ClientCpfDTO{
		ID:      1,
		Name:    "Cliente Teste",
		Email:   "teste@email.com",
		CPF:     "12345678909",
		Status:  true,
		Version: 1,
	}

	model := ToClientCpfModel(dtoInput)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, "Cliente Teste", model.Name)
	assert.Equal(t, "teste@email.com", model.Email)
	assert.Equal(t, "12345678909", model.CPF)
	assert.True(t, model.Status)
	assert.Equal(t, 1, model.Version)
}

func TestToClientCpfDTO(t *testing.T) {
	now := time.Now()

	modelInput := &models.ClientCpf{
		ID:        1,
		Name:      "Cliente DTO",
		Email:     "dto@email.com",
		CPF:       "12345678909",
		Status:    true,
		Version:   2,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dtoOutput := ToClientCpfDTO(modelInput)

	assert.Equal(t, int64(1), dtoOutput.ID)
	assert.Equal(t, "Cliente DTO", dtoOutput.Name)
	assert.Equal(t, "dto@email.com", dtoOutput.Email)
	assert.Equal(t, "12345678909", dtoOutput.CPF)
	assert.True(t, dtoOutput.Status)
	assert.Equal(t, 2, dtoOutput.Version)
	assert.NotEmpty(t, dtoOutput.CreatedAt)
	assert.NotEmpty(t, dtoOutput.UpdatedAt)

	t.Run("nil model retorna DTO vazio", func(t *testing.T) {
		dto := ToClientCpfDTO(nil)
		assert.Equal(t, ClientCpfDTO{}, dto)
	})
}

func TestToClientCpfDTOs(t *testing.T) {
	client1 := &models.ClientCpf{
		ID:     1,
		Name:   "Cliente 1",
		Email:  "c1@email.com",
		CPF:    "12345678909",
		Status: true,
	}

	client2 := &models.ClientCpf{
		ID:     2,
		Name:   "Cliente 2",
		Email:  "c2@email.com",
		CPF:    "98765432100",
		Status: false,
	}

	t.Run("lista com múltiplos clientes", func(t *testing.T) {
		input := []*models.ClientCpf{client1, client2}
		result := ToClientCpfDTOs(input)

		assert.Len(t, result, 2)
		assert.Equal(t, client1.ID, result[0].ID)
		assert.Equal(t, client2.ID, result[1].ID)
	})

	t.Run("lista com elemento nil", func(t *testing.T) {
		input := []*models.ClientCpf{client1, nil, client2}
		result := ToClientCpfDTOs(input)

		assert.Len(t, result, 2)
		assert.Equal(t, client1.ID, result[0].ID)
		assert.Equal(t, client2.ID, result[1].ID)
	})

	t.Run("lista vazia", func(t *testing.T) {
		input := []*models.ClientCpf{}
		result := ToClientCpfDTOs(input)

		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}
