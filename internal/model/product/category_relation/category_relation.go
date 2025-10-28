package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ProductCategoryRelation struct {
	ProductID  int64
	CategoryID int64
	CreatedAt  time.Time
}

func (ucr *ProductCategoryRelation) Validate() error {
	if ucr.ProductID <= 0 {
		return &validators.ValidationError{Field: "ProductID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	if ucr.CategoryID <= 0 {
		return &validators.ValidationError{Field: "CategoryID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	return nil
}
