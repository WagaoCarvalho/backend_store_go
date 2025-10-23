package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relation"
)

type SupplierCategoryRelationsDTO struct {
	SupplierID *int64     `json:"supplier_id,omitempty"`
	CategoryID *int64     `json:"category_id,omitempty"`
	Version    *int       `json:"version,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

// Converte DTO para Model
func ToSupplierCategoryRelationsModel(dto SupplierCategoryRelationsDTO) *models.SupplierCategoryRelation {
	c := &models.SupplierCategoryRelation{}

	if dto.SupplierID != nil {
		c.SupplierID = *dto.SupplierID
	}
	if dto.CategoryID != nil {
		c.CategoryID = *dto.CategoryID
	}
	if dto.Version != nil {
		c.Version = *dto.Version
	}
	if dto.CreatedAt != nil {
		c.CreatedAt = *dto.CreatedAt
	}

	return c
}

// Converte Model para DTO
func ToSupplierCategoryRelationsDTO(model *models.SupplierCategoryRelation) SupplierCategoryRelationsDTO {
	return SupplierCategoryRelationsDTO{
		SupplierID: &model.SupplierID,
		CategoryID: &model.CategoryID,
		Version:    &model.Version,
		CreatedAt:  &model.CreatedAt,
	}
}

func ToSupplierRelatiosDTOs(models []*models.SupplierCategoryRelation) []SupplierCategoryRelationsDTO {
	if len(models) == 0 {
		return []SupplierCategoryRelationsDTO{}
	}

	dtos := make([]SupplierCategoryRelationsDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToSupplierCategoryRelationsDTO(m))
		}
	}
	return dtos
}
