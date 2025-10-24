package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type UserCategoryRelation struct {
	UserID     int64
	CategoryID int64
	CreatedAt  time.Time
}

func (ucr *UserCategoryRelation) Validate() error {
	if ucr.UserID <= 0 {
		return &validators.ValidationError{Field: "UserID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	if ucr.CategoryID <= 0 {
		return &validators.ValidationError{Field: "CategoryID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	return nil
}
