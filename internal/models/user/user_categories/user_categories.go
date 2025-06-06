package models

import (
	"strings"
	"time"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

type UserCategory struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (uc *UserCategory) Validate() error {
	if utils_validators.IsBlank(uc.Name) {
		return &utils_errors.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}

	if len(uc.Name) > 100 {
		return &utils_errors.ValidationError{Field: "Name", Message: "máximo de 100 caracteres"}
	}

	if len(strings.TrimSpace(uc.Description)) > 255 {
		return &utils_errors.ValidationError{Field: "Description", Message: "máximo de 255 caracteres"}
	}

	return nil
}
