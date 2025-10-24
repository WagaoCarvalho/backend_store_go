package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type SupplierCategoryRelation struct {
	SupplierID int64
	CategoryID int64
	Version    int
	CreatedAt  time.Time
}

func (r *SupplierCategoryRelation) Validate() error {
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
