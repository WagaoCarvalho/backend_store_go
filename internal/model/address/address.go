package model

import (
	"time"

	valAddress "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/address"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Address struct {
	ID           int64
	UserID       *int64
	ClientID     *int64
	SupplierID   *int64
	Street       string
	StreetNumber string
	Complement   string
	City         string
	State        string
	Country      string
	PostalCode   string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (a *Address) Validate() error {
	var errs validators.ValidationErrors

	// --- Associação ---
	if !validators.ValidateSingleNonNil(a.UserID, a.ClientID, a.SupplierID) {
		errs = append(errs, validators.ValidationError{
			Field:   "user_id/client_id/supplier_id",
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

	// --- StreetNumber ---
	if len(a.StreetNumber) > 20 {
		errs = append(errs, validators.ValidationError{Field: "street_number", Message: "street_number max 20 characters"})
	}
	if validators.IsBlank(a.StreetNumber) {
		errs = append(errs, validators.ValidationError{
			Field:   "street_number",
			Message: "street_number é obrigatório",
		})
	}

	// --- Complement ---
	if len(a.Complement) > 255 {
		errs = append(errs, validators.ValidationError{Field: "complement", Message: "complement max 255 characters"})
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
	} else if !valAddress.IsValidBrazilianState(a.State) {
		errs = append(errs, validators.ValidationError{Field: "state", Message: validators.MsgInvalidState})
	}

	// --- Country ---
	if validators.IsBlank(a.Country) {
		errs = append(errs, validators.ValidationError{Field: "country", Message: validators.MsgRequiredField})
	} else if !validators.EqualsIgnoreCaseAndTrim(a.Country, "Brasil") {
		errs = append(errs, validators.ValidationError{Field: "country", Message: validators.MsgInvalidCountry})
	}

	// --- PostalCode ---
	if validators.IsBlank(a.PostalCode) {
		errs = append(errs, validators.ValidationError{Field: "postal_code", Message: validators.MsgRequiredField})
	} else if !valAddress.IsValidPostalCode(a.PostalCode) {
		errs = append(errs, validators.ValidationError{Field: "postal_code", Message: validators.MsgInvalidPostalCode})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
