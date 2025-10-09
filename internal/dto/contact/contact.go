package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
)

type ContactDTO struct {
	ID                 *int64     `json:"id,omitempty"`
	ContactName        string     `json:"contact_name"`
	ContactDescription string     `json:"contact_description,omitempty"`
	Email              string     `json:"email,omitempty"`
	Phone              string     `json:"phone,omitempty"`
	Cell               string     `json:"cell,omitempty"`
	ContactType        string     `json:"contact_type,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

func ToContactModel(dto ContactDTO) *models.Contact {
	c := &models.Contact{
		ContactName:        dto.ContactName,
		ContactDescription: dto.ContactDescription,
		Email:              dto.Email,
		Phone:              dto.Phone,
		Cell:               dto.Cell,
		ContactType:        dto.ContactType,
	}

	if dto.ID != nil {
		c.ID = *dto.ID
	}
	if dto.CreatedAt != nil {
		c.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		c.UpdatedAt = *dto.UpdatedAt
	}

	return c
}

func ToContactDTO(model *models.Contact) ContactDTO {
	return ContactDTO{
		ID:                 &model.ID,
		ContactName:        model.ContactName,
		ContactDescription: model.ContactDescription,
		Email:              model.Email,
		Phone:              model.Phone,
		Cell:               model.Cell,
		ContactType:        model.ContactType,
		CreatedAt:          &model.CreatedAt,
		UpdatedAt:          &model.UpdatedAt,
	}
}

func ToContactDTOs(models []*models.Contact) []ContactDTO {
	if len(models) == 0 {
		return []ContactDTO{}
	}

	dtos := make([]ContactDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToContactDTO(m))
		}
	}
	return dtos
}
