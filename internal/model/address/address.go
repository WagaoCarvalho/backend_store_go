package models

import (
	"strings"
	"time"

	val_address "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/address"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
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
	var errs validators.ValidationErrors

	// --- Associação ---
	if !validators.ValidateSingleNonNil(a.UserID, a.ClientID, a.SupplierID) {
		errs = append(errs, validators.ValidationError{
			Field:   "UserID/ClientID/SupplierID",
			Message: validators.MsgInvalidAssociation,
		})
	}

	// --- Street ---
	if validators.IsBlank(a.Street) {
		errs = append(errs, validators.ValidationError{Field: "street", Message: validators.MsgRequiredField})
	} else {
		if len(a.Street) < 3 {
			errs = append(errs, validators.ValidationError{Field: "street", Message: validators.MsgMin3})
		}
		if len(a.Street) > 100 {
			errs = append(errs, validators.ValidationError{Field: "street", Message: validators.MsgMax100})
		}
	}

	// --- City ---
	if validators.IsBlank(a.City) {
		errs = append(errs, validators.ValidationError{Field: "city", Message: validators.MsgRequiredField})
	} else if len(a.City) < 2 {
		errs = append(errs, validators.ValidationError{Field: "city", Message: validators.MsgMin2})
	}

	// --- State ---
	if validators.IsBlank(a.State) {
		errs = append(errs, validators.ValidationError{Field: "state", Message: validators.MsgRequiredField})
	} else if !val_address.IsValidBrazilianState(a.State) {
		errs = append(errs, validators.ValidationError{Field: "state", Message: validators.MsgInvalidState})
	}

	// --- Country ---
	if validators.IsBlank(a.Country) {
		errs = append(errs, validators.ValidationError{Field: "country", Message: validators.MsgRequiredField})
	} else if !strings.EqualFold(strings.TrimSpace(a.Country), "Brasil") {
		errs = append(errs, validators.ValidationError{Field: "country", Message: validators.MsgInvalidCountry})
	}

	// --- PostalCode ---
	if validators.IsBlank(a.PostalCode) {
		errs = append(errs, validators.ValidationError{Field: "postal_code", Message: validators.MsgRequiredField})
	} else if !val_address.IsValidPostalCode(a.PostalCode) {
		errs = append(errs, validators.ValidationError{Field: "postal_code", Message: validators.MsgInvalidPostalCode})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
