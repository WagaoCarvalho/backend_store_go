package models

import "time"

type SupplierCategoryRelations struct {
	ID         int64     `json:"id"`
	SupplierID int64     `json:"user_id"`
	CategoryID int64     `json:"category_id"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
