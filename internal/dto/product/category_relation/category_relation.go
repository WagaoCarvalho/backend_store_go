package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
)

type ProductCategoryRelationsDTO struct {
	ProductID  int64  `json:"product_id"`
	CategoryID int64  `json:"category_id"`
	CreatedAt  string `json:"created_at,omitempty"`
}

func ToProductCategoryRelationsModel(dto ProductCategoryRelationsDTO) *models.ProductCategoryRelation {
	return &models.ProductCategoryRelation{
		ProductID:  dto.ProductID,
		CategoryID: dto.CategoryID,
		CreatedAt:  time.Now(),
	}
}

func ToProductCategoryRelationsDTO(m *models.ProductCategoryRelation) ProductCategoryRelationsDTO {
	if m == nil {
		return ProductCategoryRelationsDTO{}
	}

	return ProductCategoryRelationsDTO{
		ProductID:  m.ProductID,
		CategoryID: m.CategoryID,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	}
}
