package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
)

type ProductCategoryRelationDTO struct {
	ProductID  int64  `json:"product_id"`
	CategoryID int64  `json:"category_id"`
	CreatedAt  string `json:"created_at,omitempty"`
}

func ToModel(dto ProductCategoryRelationDTO) *models.ProductCategoryRelation {
	return &models.ProductCategoryRelation{
		ProductID:  dto.ProductID,
		CategoryID: dto.CategoryID,
	}
}

func ToDTO(m *models.ProductCategoryRelation) ProductCategoryRelationDTO {
	if m == nil {
		return ProductCategoryRelationDTO{}
	}

	dto := ProductCategoryRelationDTO{
		ProductID:  m.ProductID,
		CategoryID: m.CategoryID,
	}

	if !m.CreatedAt.IsZero() {
		dto.CreatedAt = m.CreatedAt.Format(time.RFC3339)
	}

	return dto
}

func ToDTOs(models []*models.ProductCategoryRelation) []ProductCategoryRelationDTO {
	if len(models) == 0 {
		return []ProductCategoryRelationDTO{}
	}

	dtos := make([]ProductCategoryRelationDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToDTO(m))
		}
	}
	return dtos
}
