package models

import "time"

type SupplierCategoryRelations struct {
	SupplierID int64     `json:"supplier_id"`
	CategoryID int64     `json:"category_id"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
