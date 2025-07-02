package models

import (
	"strings"
	"time"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
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

	if a.UserID == nil && a.ClientID == nil && a.SupplierID == nil {
		return &utils_errors.ValidationError{
			Field:   "UserID/ClientID/SupplierID",
			Message: "pelo menos um deve ser informado",
		}
	}

	// --- Street ---
	if utils_validators.IsBlank(a.Street) {
		return &utils_errors.ValidationError{Field: "Street", Message: "campo obrigatório"}
	}
	if len(a.Street) < 3 {
		return &utils_errors.ValidationError{Field: "Street", Message: "mínimo de 3 caracteres"}
	}
	if len(a.Street) > 100 {
		return &utils_errors.ValidationError{Field: "Street", Message: "máximo de 100 caracteres"}
	}

	// --- City ---
	if utils_validators.IsBlank(a.City) {
		return &utils_errors.ValidationError{Field: "City", Message: "campo obrigatório"}
	}
	if len(a.City) < 2 {
		return &utils_errors.ValidationError{Field: "City", Message: "mínimo de 2 caracteres"}
	}

	// --- State ---
	if utils_validators.IsBlank(a.State) {
		return &utils_errors.ValidationError{Field: "State", Message: "campo obrigatório"}
	}
	if !utils_validators.IsValidBrazilianState(a.State) {
		return &utils_errors.ValidationError{Field: "State", Message: "estado inválido"}
	}

	// --- Country ---
	if utils_validators.IsBlank(a.Country) {
		return &utils_errors.ValidationError{Field: "Country", Message: "campo obrigatório"}
	}
	if strings.ToLower(strings.TrimSpace(a.Country)) != "brasil" {
		return &utils_errors.ValidationError{Field: "Country", Message: "país não suportado"}
	}

	// --- PostalCode ---
	if utils_validators.IsBlank(a.PostalCode) {
		return &utils_errors.ValidationError{Field: "PostalCode", Message: "campo obrigatório"}
	}
	if !utils_validators.IsValidPostalCode(a.PostalCode) {
		return &utils_errors.ValidationError{Field: "PostalCode", Message: "formato inválido (ex: 00000-000)"}
	}

	return nil
}
