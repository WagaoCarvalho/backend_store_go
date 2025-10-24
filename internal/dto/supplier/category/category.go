package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
)

type SupplierCategoryDTO struct {
	ID          *int64     `json:"id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func ToSupplierCategoryModel(dto SupplierCategoryDTO) *models.SupplierCategory {
	c := &models.SupplierCategory{
		Name:        dto.Name,
		Description: dto.Description,
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

func ToSupplierCategoryDTO(model *models.SupplierCategory) SupplierCategoryDTO {
	return SupplierCategoryDTO{
		ID:          &model.ID,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   &model.CreatedAt,
		UpdatedAt:   &model.UpdatedAt,
	}
}

func ToSupplierCategoryDTOs(models []*models.SupplierCategory) []SupplierCategoryDTO {
	if len(models) == 0 {
		return []SupplierCategoryDTO{}
	}

	dtos := make([]SupplierCategoryDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToSupplierCategoryDTO(m))
		}
	}
	return dtos
}
