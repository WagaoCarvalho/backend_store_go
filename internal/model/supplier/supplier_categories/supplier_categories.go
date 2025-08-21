package models

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
)

type SupplierCategory struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (sc *SupplierCategory) Validate() error {
	if validators.IsBlank(sc.Name) {
		return &validators.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}
	if len(sc.Name) > 100 {
		return &validators.ValidationError{Field: "Name", Message: "máximo de 100 caracteres"}
	}
	if len(sc.Description) > 255 {
		return &validators.ValidationError{Field: "Description", Message: "máximo de 255 caracteres"}
	}
	return nil
}
