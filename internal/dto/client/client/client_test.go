package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	client "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func TestToClientModel(t *testing.T) {
	id := int64(1)
	dtoInput := ClientDTO{
		ID:         &id,
		Name:       "Cliente Teste",
		Email:      utils.StrToPtr("teste@email.com"),
		CPF:        utils.StrToPtr("12345678909"),
		CNPJ:       nil,
		ClientType: "PF",
		Status:     true,
		Version:    1,
	}

	model := ToClientModel(dtoInput)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, "Cliente Teste", model.Name)
	assert.Equal(t, "teste@email.com", *model.Email)
	assert.Equal(t, "12345678909", *model.CPF)
	assert.Nil(t, model.CNPJ)
	assert.Equal(t, "PF", model.ClientType)
	assert.True(t, model.Status)
	assert.Equal(t, 1, model.Version)
}

func TestToClientDTO(t *testing.T) {
	modelInput := &client.Client{
		ID:         1,
		Name:       "Cliente DTO",
		Email:      utils.StrToPtr("dto@email.com"),
		CPF:        utils.StrToPtr("12345678909"),
		CNPJ:       nil,
		ClientType: "PF",
		Status:     true,
		Version:    2,
	}

	dtoOutput := ToClientDTO(modelInput)

	assert.NotNil(t, dtoOutput.ID)
	assert.Equal(t, int64(1), *dtoOutput.ID)
	assert.Equal(t, "Cliente DTO", dtoOutput.Name)
	assert.Equal(t, "dto@email.com", *dtoOutput.Email)
	assert.Equal(t, "12345678909", *dtoOutput.CPF)
	assert.Nil(t, dtoOutput.CNPJ)
	assert.Equal(t, "PF", dtoOutput.ClientType)
	assert.True(t, dtoOutput.Status)
	assert.Equal(t, 2, dtoOutput.Version)

	t.Run("ToClientDTO with nil input returns empty DTO", func(t *testing.T) {
		dto := ToClientDTO(nil)

		assert.Equal(t, ClientDTO{}, dto)
	})

}
