package models

import (
	"time"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

type SupplierCategory struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (sc *SupplierCategory) Validate() error {
	if utils_validators.IsBlank(sc.Name) {
		return &utils_errors.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}
	if len(sc.Name) > 100 {
		return &utils_errors.ValidationError{Field: "Name", Message: "máximo de 100 caracteres"}
	}
	if len(sc.Description) > 255 {
		return &utils_errors.ValidationError{Field: "Description", Message: "máximo de 255 caracteres"}
	}
	return nil
}
