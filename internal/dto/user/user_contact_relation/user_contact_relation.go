package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relations"
)

type UserContactRelationDTO struct {
	ContactID int64  `json:"contact_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at,omitempty"`
}

func ToContactRelationModel(dto UserContactRelationDTO) *models.UserContactRelations {
	return &models.UserContactRelations{
		ContactID: dto.ContactID,
		UserID:    dto.UserID,
		CreatedAt: time.Now(),
	}
}

func ToContactRelationDTO(m *models.UserContactRelations) UserContactRelationDTO {
	if m == nil {
		return UserContactRelationDTO{}
	}

	return UserContactRelationDTO{
		ContactID: m.ContactID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
}

func ToUserContactRelationsDTOs(models []*models.UserContactRelations) []UserContactRelationDTO {
	if len(models) == 0 {
		return []UserContactRelationDTO{}
	}

	dtos := make([]UserContactRelationDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToContactRelationDTO(m))
		}
	}
	return dtos
}
