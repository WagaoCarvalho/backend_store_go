package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
)

type ProductDTO struct {
	ID                 *int64     `json:"id,omitempty"`
	SupplierID         *int64     `json:"supplier_id,omitempty"`
	ProductName        string     `json:"product_name"`
	Manufacturer       string     `json:"manufacturer"`
	Description        string     `json:"description,omitempty"`
	CostPrice          float64    `json:"cost_price"`
	SalePrice          float64    `json:"sale_price"`
	StockQuantity      int        `json:"stock_quantity"`
	MinStock           int        `json:"min_stock"`
	MaxStock           *int       `json:"max_stock,omitempty"`
	Barcode            *string    `json:"barcode,omitempty"`
	Status             bool       `json:"status"`
	Version            int        `json:"version"`
	AllowDiscount      bool       `json:"allow_discount"`
	MinDiscountPercent float64    `json:"min_discount_percent"`
	MaxDiscountPercent float64    `json:"max_discount_percent"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

func ToProductModel(dto ProductDTO) *models.Product {
	model := &models.Product{
		ID:                 getOrDefault(dto.ID),
		SupplierID:         dto.SupplierID,
		ProductName:        dto.ProductName,
		Manufacturer:       dto.Manufacturer,
		Description:        dto.Description,
		CostPrice:          dto.CostPrice,
		SalePrice:          dto.SalePrice,
		StockQuantity:      dto.StockQuantity,
		MinStock:           dto.MinStock,
		MaxStock:           dto.MaxStock,
		Barcode:            dto.Barcode,
		Status:             dto.Status,
		Version:            dto.Version,
		AllowDiscount:      dto.AllowDiscount,
		MinDiscountPercent: dto.MinDiscountPercent,
		MaxDiscountPercent: dto.MaxDiscountPercent,
	}

	if dto.CreatedAt != nil {
		model.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		model.UpdatedAt = *dto.UpdatedAt
	}

	return model
}

func ToProductDTO(model *models.Product) ProductDTO {
	return ProductDTO{
		ID:                 &model.ID,
		SupplierID:         model.SupplierID,
		ProductName:        model.ProductName,
		Manufacturer:       model.Manufacturer,
		Description:        model.Description,
		CostPrice:          model.CostPrice,
		SalePrice:          model.SalePrice,
		StockQuantity:      model.StockQuantity,
		MinStock:           model.MinStock,
		MaxStock:           model.MaxStock,
		Barcode:            model.Barcode,
		Status:             model.Status,
		Version:            model.Version,
		AllowDiscount:      model.AllowDiscount,
		MinDiscountPercent: model.MinDiscountPercent,
		MaxDiscountPercent: model.MaxDiscountPercent,
		CreatedAt:          &model.CreatedAt,
		UpdatedAt:          &model.UpdatedAt,
	}
}

func getOrDefault(id *int64) int64 {
	if id == nil {
		return 0
	}
	return *id
}

func ToProductDTOs(models []*models.Product) []ProductDTO {
	if len(models) == 0 {
		return []ProductDTO{}
	}

	dtos := make([]ProductDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToProductDTO(m))
		}
	}
	return dtos
}
