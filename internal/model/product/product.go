package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Product struct {
	ID                 int64     `json:"id"`
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

var barcodeRegex = regexp.MustCompile(`^[0-9]{8,14}$`)

func (p *Product) Validate() error {
	if strings.TrimSpace(p.ProductName) == "" {
		return fmt.Errorf("%w: nome do produto inválido", ErrInvalidProductName)
	}
	if strings.TrimSpace(p.Manufacturer) == "" {
		return ErrInvalidManufacturer
	}
	if p.CostPrice < 0 {
		return ErrInvalidCostPrice
	}
	if p.SalePrice < 0 {
		return ErrInvalidSalePrice
	}
	if p.SalePrice < p.CostPrice {
		return ErrSalePriceBelowCost
	}
	if p.StockQuantity < 0 {
		return ErrNegativeStock
	}
	if p.Barcode != "" && !barcodeRegex.MatchString(p.Barcode) {
		return ErrInvalidBarcode
	}
	if p.SupplierID == nil {
		return ErrSupplierRequired
	}
	if !p.Status {
		return ErrInactiveProductNotAllowed
	}

	// Regras de desconto (apenas validação estrutural)
	if p.AllowDiscount {
		if p.MinDiscountPercent < 0 || p.MaxDiscountPercent < 0 {
			return ErrNegativeDiscount
		}
		if p.MinDiscountPercent > p.MaxDiscountPercent {
			return ErrInvalidDiscountRange
		}
		if p.MaxDiscountPercent > 100 {
			return ErrDiscountAboveLimit
		}
	}

	return nil
}
