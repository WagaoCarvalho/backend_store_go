package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ProductCategoryRelation struct {
	ProductID  int64     `json:"product_id"`
	CategoryID int64     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (pcr *ProductCategoryRelation) Validate() error {
	if pcr.ProductID <= 0 {
		return &validators.ValidationError{
			Field:   "product_id",
			Message: "ID do produto é obrigatório e deve ser maior que zero",
		}
	}

	if pcr.CategoryID <= 0 {
		return &validators.ValidationError{
			Field:   "category_id",
			Message: "ID da categoria é obrigatório e deve ser maior que zero",
		}
	}

	return nil
}
