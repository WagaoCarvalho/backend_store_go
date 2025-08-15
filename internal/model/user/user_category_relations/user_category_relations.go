package models

import (
	"time"

	err "github.com/WagaoCarvalho/backend_store_go/pkg/utils"
)

type UserCategoryRelations struct {
	UserID     int64     `json:"user_id"`
	CategoryID int64     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ucr *UserCategoryRelations) Validate() error {
	if ucr.UserID <= 0 {
		return &err.ValidationError{Field: "UserID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	if ucr.CategoryID <= 0 {
		return &err.ValidationError{Field: "CategoryID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	return nil
}
