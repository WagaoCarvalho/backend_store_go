package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
)

type ClientCpfDTO struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CPF         string `json:"cpf"`
	Description string `json:"description,omitempty"`
	Version     int    `json:"version"`
	Status      bool   `json:"status"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func ToClientCpfModel(dto ClientCpfDTO) *models.ClientCpf {
	return &models.ClientCpf{
		ID:          dto.ID,
		Name:        dto.Name,
		Email:       dto.Email,
		CPF:         dto.CPF,
		Description: dto.Description,
		Version:     dto.Version,
		Status:      dto.Status,
	}
}

func ToClientCpfDTO(m *models.ClientCpf) ClientCpfDTO {
	if m == nil {
		return ClientCpfDTO{}
	}

	return ClientCpfDTO{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		CPF:         m.CPF,
		Description: m.Description,
		Version:     m.Version,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
}

func ToClientCpfDTOs(models []*models.ClientCpf) []ClientCpfDTO {
	if len(models) == 0 {
		return []ClientCpfDTO{}
	}

	dtos := make([]ClientCpfDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToClientCpfDTO(m))
		}
	}
	return dtos
}
