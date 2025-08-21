package models

import (
	"strings"
	"time"

	err "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
)

type UserCategory struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (uc *UserCategory) Validate() error {
	if validators.IsBlank(uc.Name) {
		return &err.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}

	if len(uc.Name) > 100 {
		return &err.ValidationError{Field: "Name", Message: "máximo de 100 caracteres"}
	}

	if len(strings.TrimSpace(uc.Description)) > 255 {
		return &err.ValidationError{Field: "Description", Message: "máximo de 255 caracteres"}
	}

	return nil
}
