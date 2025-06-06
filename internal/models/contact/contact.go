package models

import (
	"strings"
	"time"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
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
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (c *Contact) Validate() error {
	if c.UserID == nil && c.ClientID == nil && c.SupplierID == nil {
		return &utils_errors.ValidationError{
			Field:   "UserID/ClientID/SupplierID",
			Message: "pelo menos um deve ser informado",
		}
	}

	// --- ContactName ---
	if utils_validators.IsBlank(c.ContactName) {
		return &utils_errors.ValidationError{Field: "ContactName", Message: "campo obrigatório"}
	}
	if len(c.ContactName) < 3 {
		return &utils_errors.ValidationError{Field: "ContactName", Message: "mínimo de 3 caracteres"}
	}
	if len(c.ContactName) > 100 {
		return &utils_errors.ValidationError{Field: "ContactName", Message: "máximo de 100 caracteres"}
	}

	// --- ContactPosition ---
	if len(c.ContactPosition) > 100 {
		return &utils_errors.ValidationError{Field: "ContactPosition", Message: "máximo de 100 caracteres"}
	}

	// --- Email ---
	if !utils_validators.IsBlank(c.Email) {
		if !utils_validators.IsValidEmail(c.Email) {
			return &utils_errors.ValidationError{Field: "Email", Message: "formato inválido"}
		}
		if len(c.Email) > 100 {
			return &utils_errors.ValidationError{Field: "Email", Message: "máximo de 100 caracteres"}
		}
	}

	// --- Phone ---
	if !utils_validators.IsBlank(c.Phone) {
		if !utils_validators.IsValidPhone(c.Phone) {
			return &utils_errors.ValidationError{Field: "Phone", Message: "formato inválido (ex: (11) 1234-5678)"}
		}
	}

	// --- Cell ---
	if !utils_validators.IsBlank(c.Cell) {
		if !utils_validators.IsValidCell(c.Cell) {
			return &utils_errors.ValidationError{Field: "Cell", Message: "formato inválido (ex: (11) 91234-5678)"}
		}
	}

	// --- ContactType (opcional, mas se informado, validar) ---
	if !utils_validators.IsBlank(c.ContactType) {
		validTypes := map[string]bool{
			"principal":  true,
			"financeiro": true,
			"comercial":  true,
			"suporte":    true,
		}
		normalized := strings.ToLower(strings.TrimSpace(c.ContactType))
		if !validTypes[normalized] {
			return &utils_errors.ValidationError{Field: "ContactType", Message: "tipo inválido"}
		}
	}

	return nil
}
