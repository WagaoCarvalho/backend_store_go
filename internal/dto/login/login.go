package dto

import models "github.com/WagaoCarvalho/backend_store_go/internal/model/login"

// --- Entrada (request) ---
type LoginCredentialsDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *LoginCredentialsDTO) ToModel() *models.LoginCredential {
	if dto == nil {
		return nil
	}
	return &models.LoginCredential{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

// --- Saída (response) ---
type AuthResponseDTO struct {
	AccessToken string `json:"access_token"`
	//RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn int64  `json:"expires_in"` // em segundos
	TokenType string `json:"token_type"` // geralmente "Bearer"
}

func ToAuthResponseDTO(m *models.AuthResponse) *AuthResponseDTO {
	if m == nil {
		return nil
	}
	return &AuthResponseDTO{
		AccessToken: m.AccessToken,
		//RefreshToken: m.RefreshToken,
		ExpiresIn: m.ExpiresIn,
		TokenType: m.TokenType,
	}
}
