package model

import (
	valContact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type LoginCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string
	//RefreshToken string
	ExpiresIn int64
	TokenType string
}

func (c *LoginCredential) Validate() error {
	var errs validators.ValidationErrors

	// Email
	if validators.IsBlank(c.Email) {
		errs = append(errs, validators.ValidationError{
			Field:   "email",
			Message: validators.MsgRequiredField,
		})
	} else if len(c.Email) > 100 {
		errs = append(errs, validators.ValidationError{
			Field:   "email",
			Message: validators.MsgMax100,
		})
	} else if !valContact.IsValidEmail(c.Email) {
		errs = append(errs, validators.ValidationError{
			Field:   "email",
			Message: validators.MsgInvalidEmail,
		})
	}

	// Password: apenas verificar preenchimento
	if validators.IsBlank(c.Password) {
		errs = append(errs, validators.ValidationError{Field: "password", Message: validators.MsgRequiredField})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
