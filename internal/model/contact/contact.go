package model

import (
	"strings"
	"time"

	valContact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Contact struct {
	ID                 int64
	ContactName        string
	ContactDescription string
	Email              string
	Phone              string
	Cell               string
	ContactType        string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (c *Contact) Validate() error {
	var errs validators.ValidationErrors

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

	// ContactDescription
	if len(c.ContactDescription) > 100 {
		errs = append(errs, validators.ValidationError{Field: "contact_description", Message: validators.MsgMax100})
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
