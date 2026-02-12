package model

import (
	"strings"
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ProductCategory struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (pc *ProductCategory) Validate() error {
	var errs validators.ValidationErrors

	// --- Name ---
	name := strings.TrimSpace(pc.Name)
	if validators.IsBlank(name) {
		errs = append(errs, validators.ValidationError{
			Field:   "name",
			Message: validators.MsgRequiredField,
		})
	} else {
		if len(name) < 2 {
			errs = append(errs, validators.ValidationError{
				Field:   "name",
				Message: validators.MsgMin2,
			})
		}
		if len(name) > 255 {
			errs = append(errs, validators.ValidationError{
				Field:   "name",
				Message: "nome máximo 255 caracteres",
			})
		}
	}

	// --- Description ---
	description := strings.TrimSpace(pc.Description)
	if description != "" && len(description) > 255 {
		errs = append(errs, validators.ValidationError{
			Field:   "description",
			Message: "descrição máxima 255 caracteres",
		})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
