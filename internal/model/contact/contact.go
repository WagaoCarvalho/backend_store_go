package models

import (
	"strings"
	"time"

	err "github.com/WagaoCarvalho/backend_store_go/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/pkg/utils/validators"
)

type Contact struct {
	ID              int64     `json:"id"`
	UserID          *int64    `json:"user_id,omitempty"`
	ClientID        *int64    `json:"client_id,omitempty"`
	SupplierID      *int64    `json:"supplier_id,omitempty"`
	ContactName     string    `json:"contact_name"`
	ContactPosition string    `json:"contact_position,omitempty"`
	Email           string    `json:"email,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	Cell            string    `json:"cell,omitempty"`
	ContactType     string    `json:"contact_type,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (c *Contact) Validate() error {
	if !validators.ValidateSingleNonNil(c.UserID, c.ClientID, c.SupplierID) {
		return &err.ValidationError{
			Field:   "UserID/ClientID/SupplierID",
			Message: "exatamente um deve ser informado",
		}
	}

	if validators.IsBlank(c.ContactName) {
		return &err.ValidationError{Field: "ContactName", Message: "campo obrigatório"}
	}
	if len(c.ContactName) < 3 {
		return &err.ValidationError{Field: "ContactName", Message: "mínimo de 3 caracteres"}
	}
	if len(c.ContactName) > 100 {
		return &err.ValidationError{Field: "ContactName", Message: "máximo de 100 caracteres"}
	}

	if len(c.ContactPosition) > 100 {
		return &err.ValidationError{Field: "ContactPosition", Message: "máximo de 100 caracteres"}
	}

	if !validators.IsBlank(c.Email) {
		if !validators.IsValidEmail(c.Email) {
			return &err.ValidationError{Field: "Email", Message: "formato inválido"}
		}
		if len(c.Email) > 100 {
			return &err.ValidationError{Field: "Email", Message: "máximo de 100 caracteres"}
		}
	}

	if !validators.IsBlank(c.Phone) {
		if !validators.IsValidPhone(c.Phone) {
			return &err.ValidationError{Field: "Phone", Message: "formato inválido (ex: 1112345678)"}
		}
	}

	if !validators.IsBlank(c.Cell) {
		if !validators.IsValidCell(c.Cell) {
			return &err.ValidationError{Field: "Cell", Message: "formato inválido (ex: 11912345678)"}
		}
	}

	if !validators.IsBlank(c.ContactType) {
		validTypes := map[string]bool{
			"principal":  true,
			"financeiro": true,
			"comercial":  true,
			"suporte":    true,
		}
		normalized := strings.ToLower(strings.TrimSpace(c.ContactType))
		if !validTypes[normalized] {
			return &err.ValidationError{Field: "ContactType", Message: "tipo inválido"}
		}
	}

	return nil
}
