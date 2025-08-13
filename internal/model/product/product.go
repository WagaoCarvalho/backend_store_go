package models

import (
	"regexp"
	"strings"
	"time"
)

type Product struct {
	ID            int64     `json:"id"`
	SupplierID    *int64    `json:"supplier_id,omitempty"`
	ProductName   string    `json:"product_name"`
	Manufacturer  string    `json:"manufacturer"`
	Description   string    `json:"product_description,omitempty"`
	CostPrice     float64   `json:"cost_price"`
	SalePrice     float64   `json:"sale_price"`
	StockQuantity int       `json:"stock_quantity"`
	Barcode       string    `json:"barcode,omitempty"`
	Status        bool      `json:"status"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var barcodeRegex = regexp.MustCompile(`^[0-9]{8,14}$`)

func (p *Product) Validate() error {
	if strings.TrimSpace(p.ProductName) == "" {
		return ErrInvalidProductName
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
	return nil
}
