package models

import (
	err "github.com/WagaoCarvalho/backend_store_go/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/pkg/utils/validators"
)

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *LoginCredentials) Validate() error {
	if validators.IsBlank(c.Email) {
		return &err.ValidationError{Field: "Email", Message: "campo obrigatório"}
	}

	if len(c.Email) > 100 {
		return &err.ValidationError{Field: "Email", Message: "máximo de 100 caracteres"}
	}

	if !validators.IsValidEmail(c.Email) {
		return &err.ValidationError{Field: "Email", Message: "email inválido"}
	}

	if err := validators.ValidateStrongPassword(c.Password); err != nil {
		return err
	}

	return nil
}
