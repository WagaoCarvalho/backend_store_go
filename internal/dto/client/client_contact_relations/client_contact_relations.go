package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client_contact_relations"
)

type ClientContactRelationDTO struct {
	ContactID int64  `json:"contact_id"`
	ClientID  int64  `json:"client_id"`
	CreatedAt string `json:"created_at,omitempty"`
}

func ToClientContactRelationModel(dto ClientContactRelationDTO) *models.ClientContactRelations {
	return &models.ClientContactRelations{
		ContactID: dto.ContactID,
		ClientID:  dto.ClientID,
		CreatedAt: time.Now(),
	}
}

func ToClientContactRelationDTO(m *models.ClientContactRelations) ClientContactRelationDTO {
	if m == nil {
		return ClientContactRelationDTO{}
	}

	return ClientContactRelationDTO{
		ContactID: m.ContactID,
		ClientID:  m.ClientID,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
}
