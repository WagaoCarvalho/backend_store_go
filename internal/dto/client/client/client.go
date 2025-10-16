package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type ClientDTO struct {
	ID          *int64  `json:"id,omitempty"`
	Name        string  `json:"name"`
	Email       *string `json:"email,omitempty"`
	CPF         *string `json:"cpf,omitempty"`
	CNPJ        *string `json:"cnpj,omitempty"`
	Description string  `json:"description,omitempty"`
	Version     int     `json:"version"`
	Status      bool    `json:"status"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

func ToClientModel(dto ClientDTO) *models.Client {
	return &models.Client{
		ID:          utils.NilToZero(dto.ID),
		Name:        dto.Name,
		Email:       dto.Email,
		CPF:         dto.CPF,
		CNPJ:        dto.CNPJ,
		Description: dto.Description,
		Version:     dto.Version,
		Status:      dto.Status,
	}
}

func ToClientDTO(m *models.Client) ClientDTO {
	if m == nil {
		return ClientDTO{}
	}

	return ClientDTO{
		ID:          &m.ID,
		Name:        m.Name,
		Email:       m.Email,
		CPF:         m.CPF,
		CNPJ:        m.CNPJ,
		Description: m.Description,
		Version:     m.Version,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
}

func ToClientDTOs(models []*models.Client) []ClientDTO {
	if len(models) == 0 {
		return []ClientDTO{}
	}

	dtos := make([]ClientDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToClientDTO(m))
		}
	}
	return dtos
}
