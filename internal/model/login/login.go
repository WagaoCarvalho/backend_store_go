package models

import (
	"regexp"

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

	if validators.IsBlank(c.Password) {
		return &err.ValidationError{Field: "Password", Message: "campo obrigatório"}
	}

	if len(c.Password) < 8 {
		return &err.ValidationError{Field: "Password", Message: "mínimo de 8 caracteres"}
	}

	if len(c.Password) > 64 {
		return &err.ValidationError{Field: "Password", Message: "máximo de 64 caracteres"}
	}

	// Senha forte: pelo menos uma maiúscula, uma minúscula, um número e um caractere especial
	uppercase := regexp.MustCompile(`[A-Z]`)
	lowercase := regexp.MustCompile(`[a-z]`)
	number := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':"\\|,.<>\/?]`)

	if !uppercase.MatchString(c.Password) {
		return &err.ValidationError{Field: "Password", Message: "deve conter pelo menos uma letra maiúscula"}
	}
	if !lowercase.MatchString(c.Password) {
		return &err.ValidationError{Field: "Password", Message: "deve conter pelo menos uma letra minúscula"}
	}
	if !number.MatchString(c.Password) {
		return &err.ValidationError{Field: "Password", Message: "deve conter pelo menos um número"}
	}
	if !special.MatchString(c.Password) {
		return &err.ValidationError{Field: "Password", Message: "deve conter pelo menos um caractere especial"}
	}

	return nil
}
