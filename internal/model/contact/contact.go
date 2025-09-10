package model

import (
	"strings"
	"time"

	valContact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
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
	var errs validators.ValidationErrors

	// --- Associação ---
	if !validators.ValidateSingleNonNil(c.UserID, c.ClientID, c.SupplierID) {
		errs = append(errs, validators.ValidationError{
			Field:   "UserID/ClientID/SupplierID",
			Message: validators.MsgInvalidAssociation,
		})
	}

	// ContactName
	if validators.IsBlank(c.ContactName) {
		errs = append(errs, validators.ValidationError{Field: "contact_name", Message: validators.MsgRequiredField})
	} else {
		if len(c.ContactName) < 3 {
			errs = append(errs, validators.ValidationError{Field: "contact_name", Message: validators.MsgMin3})
		}
		if len(c.ContactName) > 100 {
			errs = append(errs, validators.ValidationError{Field: "contact_name", Message: validators.MsgMax100})
		}
	}

	// ContactPosition
	if len(c.ContactPosition) > 100 {
		errs = append(errs, validators.ValidationError{Field: "contact_position", Message: validators.MsgMax100})
	}

	// Email
	if !validators.IsBlank(c.Email) {
		if !valContact.IsValidEmail(c.Email) {
			errs = append(errs, validators.ValidationError{Field: "email", Message: validators.MsgInvalidFormat})
		}
		if len(c.Email) > 100 {
			errs = append(errs, validators.ValidationError{Field: "email", Message: validators.MsgMax100})
		}
	}

	// Phone
	if !validators.IsBlank(c.Phone) {
		if !valContact.IsValidPhone(c.Phone) {
			errs = append(errs, validators.ValidationError{Field: "phone", Message: validators.MsgInvalidPhone})
		}
	}

	// Cell
	if !validators.IsBlank(c.Cell) {
		if !valContact.IsValidCell(c.Cell) {
			errs = append(errs, validators.ValidationError{Field: "cell", Message: validators.MsgInvalidCell})
		}
	}

	// ContactType
	if !validators.IsBlank(c.ContactType) {
		validTypes := map[string]bool{
			"principal":  true,
			"financeiro": true,
			"comercial":  true,
			"suporte":    true,
		}
		normalized := strings.ToLower(strings.TrimSpace(c.ContactType))
		if !validTypes[normalized] {
			errs = append(errs, validators.ValidationError{Field: "contact_type", Message: validators.MsgInvalidType})
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
