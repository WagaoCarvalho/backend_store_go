package models

import "time"

type Product struct {
	ID            int       `json:"id"`
	ProductName   string    `json:"name"`
	Manufacturer  string    `json:"manufacturer"`
	Description   string    `json:"description"`
	CostPrice     float64   `json:"cost_price"`
	SalePrice     float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	Barcode       string    `json:"barcode,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
