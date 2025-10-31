package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type SaleItemDTO struct {
	ID          *int64  `json:"id,omitempty"`
	SaleID      int64   `json:"sale_id"`
	ProductID   int64   `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Discount    float64 `json:"discount,omitempty"`
	Tax         float64 `json:"tax,omitempty"`
	Subtotal    float64 `json:"subtotal"`
	Description string  `json:"description,omitempty"`
	CreatedAt   *string `json:"created_at,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
}

// --- Conversões DTO ↔ Model ---

func ToSaleItemModel(dto SaleItemDTO) *models.SaleItem {
	model := &models.SaleItem{
		ID:          utils.NilToZero(dto.ID),
		SaleID:      dto.SaleID,
		ProductID:   dto.ProductID,
		Quantity:    dto.Quantity,
		UnitPrice:   dto.UnitPrice,
		Discount:    dto.Discount,
		Tax:         dto.Tax,
		Subtotal:    dto.Subtotal,
		Description: dto.Description,
	}

	if dto.CreatedAt != nil {
		if t, err := time.Parse(time.RFC3339, *dto.CreatedAt); err == nil {
			model.CreatedAt = t
		}
	}
	if dto.UpdatedAt != nil {
		if t, err := time.Parse(time.RFC3339, *dto.UpdatedAt); err == nil {
			model.UpdatedAt = t
		}
	}

	return model
}

func ToSaleItemDTO(model *models.SaleItem) SaleItemDTO {
	createdAt := model.CreatedAt.Format(time.RFC3339)
	updatedAt := model.UpdatedAt.Format(time.RFC3339)

	return SaleItemDTO{
		ID:          &model.ID,
		SaleID:      model.SaleID,
		ProductID:   model.ProductID,
		Quantity:    model.Quantity,
		UnitPrice:   model.UnitPrice,
		Discount:    model.Discount,
		Tax:         model.Tax,
		Subtotal:    model.Subtotal,
		Description: model.Description,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}
}

func ToSaleItemDTOList(modelsList []*models.SaleItem) []SaleItemDTO {
	result := make([]SaleItemDTO, len(modelsList))
	for i, m := range modelsList {
		result[i] = ToSaleItemDTO(m)
	}
	return result
}

func SaleItemDTOListToModelList(dtos []*SaleItemDTO) []*models.SaleItem {
	result := make([]*models.SaleItem, len(dtos))
	for i, dto := range dtos {
		result[i] = ToSaleItemModel(*dto)
	}
	return result
}

type ItemExistsResponseDTO struct {
	Exists bool `json:"exists"`
}
