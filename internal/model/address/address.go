package models

import (
	"strings"
	"time"

	err "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
)

type Address struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id,omitempty"`
	ClientID   *int64    `json:"client_id,omitempty"`
	SupplierID *int64    `json:"supplier_id,omitempty"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (a *Address) Validate() error {

	if !validators.ValidateSingleNonNil(a.UserID, a.ClientID, a.SupplierID) {
		return &err.ValidationError{
			Field:   "UserID/ClientID/SupplierID",
			Message: "exatamente um deve ser informado",
		}
	}

	// --- Street ---
	if validators.IsBlank(a.Street) {
		return &err.ValidationError{Field: "Street", Message: "campo obrigatório"}
	}
	if len(a.Street) < 3 {
		return &err.ValidationError{Field: "Street", Message: "mínimo de 3 caracteres"}
	}
	if len(a.Street) > 100 {
		return &err.ValidationError{Field: "Street", Message: "máximo de 100 caracteres"}
	}

	// --- City ---
	if validators.IsBlank(a.City) {
		return &err.ValidationError{Field: "City", Message: "campo obrigatório"}
	}
	if len(a.City) < 2 {
		return &err.ValidationError{Field: "City", Message: "mínimo de 2 caracteres"}
	}

	// --- State ---
	if validators.IsBlank(a.State) {
		return &err.ValidationError{Field: "State", Message: "campo obrigatório"}
	}
	if !validators.IsValidBrazilianState(a.State) {
		return &err.ValidationError{Field: "State", Message: "estado inválido"}
	}

	// --- Country ---
	if validators.IsBlank(a.Country) {
		return &err.ValidationError{Field: "Country", Message: "campo obrigatório"}
	}
	if strings.ToLower(strings.TrimSpace(a.Country)) != "brasil" {
		return &err.ValidationError{Field: "Country", Message: "país não suportado"}
	}

	// --- PostalCode ---
	if validators.IsBlank(a.PostalCode) {
		return &err.ValidationError{Field: "PostalCode", Message: "campo obrigatório"}
	}
	if !validators.IsValidPostalCode(a.PostalCode) {
		return &err.ValidationError{Field: "PostalCode", Message: "formato inválido (ex: 12345678)"}
	}

	return nil
}
