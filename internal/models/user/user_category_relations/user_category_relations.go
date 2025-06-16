package models

import (
	"time"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

type UserCategoryRelations struct {
	UserID     int64     `json:"user_id"`
	CategoryID int64     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ucr *UserCategoryRelations) Validate() error {
	if ucr.UserID <= 0 {
		return &utils_errors.ValidationError{Field: "UserID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	if ucr.CategoryID <= 0 {
		return &utils_errors.ValidationError{Field: "CategoryID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	return nil
}
