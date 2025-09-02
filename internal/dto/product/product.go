package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
)

type ProductDTO struct {
	ID                 *int64    `json:"id,omitempty"`
	SupplierID         *int64    `json:"supplier_id,omitempty"`
	ProductName        string    `json:"product_name"`
	Manufacturer       string    `json:"manufacturer"`
	Description        string    `json:"product_description,omitempty"`
	CostPrice          float64   `json:"cost_price"`
	SalePrice          float64   `json:"sale_price"`
	StockQuantity      int       `json:"stock_quantity"`
	Barcode            string    `json:"barcode,omitempty"`
	Status             bool      `json:"status"`
	Version            int       `json:"version"`
	AllowDiscount      bool      `json:"allow_discount"`
	MinDiscountPercent float64   `json:"min_discount_percent"`
	MaxDiscountPercent float64   `json:"max_discount_percent"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func ToProductModel(dto ProductDTO) *models.Product {
	return &models.Product{
		ID:                 getOrDefault(dto.ID),
		SupplierID:         dto.SupplierID,
		ProductName:        dto.ProductName,
		Manufacturer:       dto.Manufacturer,
		Description:        dto.Description,
		CostPrice:          dto.CostPrice,
		SalePrice:          dto.SalePrice,
		StockQuantity:      dto.StockQuantity,
		Barcode:            dto.Barcode,
		Status:             dto.Status,
		Version:            dto.Version,
		AllowDiscount:      dto.AllowDiscount,
		MinDiscountPercent: dto.MinDiscountPercent,
		MaxDiscountPercent: dto.MaxDiscountPercent,
		CreatedAt:          dto.CreatedAt,
		UpdatedAt:          dto.UpdatedAt,
	}
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
		Barcode:            model.Barcode,
		Status:             model.Status,
		Version:            model.Version,
		AllowDiscount:      model.AllowDiscount,
		MinDiscountPercent: model.MinDiscountPercent,
		MaxDiscountPercent: model.MaxDiscountPercent,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
}

func getOrDefault(id *int64) int64 {
	if id == nil {
		return 0
	}
	return *id
}
