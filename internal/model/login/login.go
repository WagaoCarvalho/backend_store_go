package models

import (
	valContact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
	valPass "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/password"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

// variável de função que pode ser sobrescrita em testes
var validateStrongPassword = valPass.ValidateStrongPassword

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *LoginCredentials) Validate() error {
	var errs validators.ValidationErrors

	// Email
	if validators.IsBlank(c.Email) {
		errs = append(errs, validators.ValidationError{Field: "email", Message: validators.MsgRequiredField})
	} else if len(c.Email) > 100 {
		errs = append(errs, validators.ValidationError{Field: "email", Message: validators.MsgMax100})
	} else if !valContact.IsValidEmail(c.Email) {
		errs = append(errs, validators.ValidationError{Field: "email", Message: validators.MsgInvalidEmail})
	}

	// Password
	if validators.IsBlank(c.Password) {
		errs = append(errs, validators.ValidationError{Field: "password", Message: validators.MsgRequiredField})
	} else {
		if err := validateStrongPassword(c.Password); err != nil {
			// Presume que ValidateStrongPassword retorna ValidationError
			if vErr, ok := err.(validators.ValidationError); ok {
				errs = append(errs, vErr)
			} else {
				// Se retornar erro genérico, adiciona com campo password
				errs = append(errs, validators.ValidationError{Field: "password", Message: err.Error()})
			}
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
