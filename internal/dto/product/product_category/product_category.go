package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product_category"
)

type ProductCategoryDTO struct {
	ID          *uint  `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// Converte DTO para Model
func ToProductCategoryModel(dto ProductCategoryDTO) *models.ProductCategory {
	id := uint(0)
	if dto.ID != nil {
		id = *dto.ID
	}

	return &models.ProductCategory{
		ID:          id,
		Name:        dto.Name,
		Description: dto.Description,
	}
}

// Converte Model para DTO
func ToProductCategoryDTO(m *models.ProductCategory) ProductCategoryDTO {
	if m == nil {
		return ProductCategoryDTO{}
	}

	return ProductCategoryDTO{
		ID:          &m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
}

func ToProductCategoryDTOs(models []*models.ProductCategory) []ProductCategoryDTO {
	if len(models) == 0 {
		return []ProductCategoryDTO{}
	}

	dtos := make([]ProductCategoryDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToProductCategoryDTO(m))
		}
	}
	return dtos
}
