package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
)

type SupplierCategoryRelationsDTO struct {
	SupplierID *int64     `json:"supplier_id,omitempty"`
	CategoryID *int64     `json:"category_id,omitempty"`
	Version    *int       `json:"version,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

// Converte DTO para Model
func ToSupplierCategoryRelationsModel(dto SupplierCategoryRelationsDTO) *models.SupplierCategoryRelations {
	c := &models.SupplierCategoryRelations{}

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
func ToSupplierCategoryRelationsDTO(model *models.SupplierCategoryRelations) SupplierCategoryRelationsDTO {
	return SupplierCategoryRelationsDTO{
		SupplierID: &model.SupplierID,
		CategoryID: &model.CategoryID,
		Version:    &model.Version,
		CreatedAt:  &model.CreatedAt,
	}
}
