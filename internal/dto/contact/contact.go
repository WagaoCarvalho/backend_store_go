package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
)

type ContactDTO struct {
	ID              *int64     `json:"id,omitempty"`
	UserID          *int64     `json:"user_id,omitempty"`
	ClientID        *int64     `json:"client_id,omitempty"`
	SupplierID      *int64     `json:"supplier_id,omitempty"`
	ContactName     string     `json:"contact_name"`
	ContactPosition string     `json:"contact_position,omitempty"`
	Email           string     `json:"email,omitempty"`
	Phone           string     `json:"phone,omitempty"`
	Cell            string     `json:"cell,omitempty"`
	ContactType     string     `json:"contact_type,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

func ToContactModel(dto ContactDTO) *models.Contact {
	c := &models.Contact{
		UserID:          dto.UserID,
		ClientID:        dto.ClientID,
		SupplierID:      dto.SupplierID,
		ContactName:     dto.ContactName,
		ContactPosition: dto.ContactPosition,
		Email:           dto.Email,
		Phone:           dto.Phone,
		Cell:            dto.Cell,
		ContactType:     dto.ContactType,
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
		ID:              &model.ID,
		UserID:          model.UserID,
		ClientID:        model.ClientID,
		SupplierID:      model.SupplierID,
		ContactName:     model.ContactName,
		ContactPosition: model.ContactPosition,
		Email:           model.Email,
		Phone:           model.Phone,
		Cell:            model.Cell,
		ContactType:     model.ContactType,
		CreatedAt:       &model.CreatedAt,
		UpdatedAt:       &model.UpdatedAt,
	}
}
