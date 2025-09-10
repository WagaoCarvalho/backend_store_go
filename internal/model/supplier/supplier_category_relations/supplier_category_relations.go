package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type SupplierCategoryRelations struct {
	SupplierID int64     `json:"supplier_id"`
	CategoryID int64     `json:"category_id"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
}

func (r *SupplierCategoryRelations) Validate() error {
	if r.SupplierID <= 0 {
		return &validators.ValidationError{
			Field:   "SupplierID",
			Message: "ID do fornecedor é obrigatório e deve ser maior que zero",
		}
	}
	if r.CategoryID <= 0 {
		return &validators.ValidationError{
			Field:   "CategoryID",
			Message: "ID da categoria é obrigatório e deve ser maior que zero",
		}
	}
	return nil
}
