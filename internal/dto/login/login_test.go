package dto

import (
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	"github.com/stretchr/testify/assert"
)

func TestLoginCredentialsDTO_ToModel(t *testing.T) {
	dto := &LoginCredentialsDTO{
		Email:    "user@example.com",
		Password: "password123",
	}

	model := dto.ToModel()

	assert.NotNil(t, model)
	assert.Equal(t, "user@example.com", model.Email)
	assert.Equal(t, "password123", model.Password)
}

func TestLoginCredentialsDTO_ToModel_Nil(t *testing.T) {
	var dto *LoginCredentialsDTO
	model := dto.ToModel()
	assert.Nil(t, model)
}

func TestToAuthResponseDTO(t *testing.T) {
	model := &models.AuthResponse{
		AccessToken: "token123",
		ExpiresIn:   3600,
		TokenType:   "Bearer",
	}

	dto := ToAuthResponseDTO(model)

	assert.NotNil(t, dto)
	assert.Equal(t, "token123", dto.AccessToken)
	assert.Equal(t, int64(3600), dto.ExpiresIn)
	assert.Equal(t, "Bearer", dto.TokenType)
}

func TestToAuthResponseDTO_Nil(t *testing.T) {
	var model *models.AuthResponse
	dto := ToAuthResponseDTO(model)
	assert.Nil(t, dto)
}
