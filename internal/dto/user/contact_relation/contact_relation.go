package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
)

type UserContactRelationDTO struct {
	ContactID int64  `json:"contact_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at,omitempty"`
}

func ToContactRelationModel(dto UserContactRelationDTO) *models.UserContactRelation {
	return &models.UserContactRelation{
		ContactID: dto.ContactID,
		UserID:    dto.UserID,
		CreatedAt: time.Now(),
	}
}

func ToContactRelationDTO(m *models.UserContactRelation) UserContactRelationDTO {
	if m == nil {
		return UserContactRelationDTO{}
	}

	return UserContactRelationDTO{
		ContactID: m.ContactID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
}

func ToUserContactRelationsDTOs(models []*models.UserContactRelation) []UserContactRelationDTO {
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
