package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_contact_relation"
)

type ContactSupplierRelationDTO struct {
	ContactID  int64  `json:"contact_id"`
	SupplierID int64  `json:"supplier_id"`
	CreatedAt  string `json:"created_at,omitempty"`
}

func ToContactSupplierRelationModel(dto ContactSupplierRelationDTO) *models.SupplierContactRelation {
	return &models.SupplierContactRelation{
		ContactID:  dto.ContactID,
		SupplierID: dto.SupplierID,
		CreatedAt:  time.Now(),
	}
}

func ToContactSupplierRelationDTO(m *models.SupplierContactRelation) ContactSupplierRelationDTO {
	if m == nil {
		return ContactSupplierRelationDTO{}
	}

	return ContactSupplierRelationDTO{
		ContactID:  m.ContactID,
		SupplierID: m.SupplierID,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	}
}

func ToSupplierContactRelationsDTOs(relations []*models.SupplierContactRelation) []ContactSupplierRelationDTO {
	if relations == nil {
		return []ContactSupplierRelationDTO{}
	}

	dtos := make([]ContactSupplierRelationDTO, len(relations))
	for i, rel := range relations {
		dtos[i] = ToContactSupplierRelationDTO(rel)
	}

	return dtos
}
