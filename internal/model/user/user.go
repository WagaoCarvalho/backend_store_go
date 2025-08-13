package models

import (
	"time"

	"unicode"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

type User struct {
	UID       int64     `json:"uid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Status    bool      `json:"status"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	// Username obrigatório e tamanho
	if utils_validators.IsBlank(u.Username) {
		return &utils_errors.ValidationError{Field: "Username", Message: "campo obrigatório"}
	}
	if len(u.Username) < 3 || len(u.Username) > 50 {
		return &utils_errors.ValidationError{Field: "Username", Message: "deve ter entre 3 e 50 caracteres"}
	}

	// Email obrigatório e válido
	if utils_validators.IsBlank(u.Email) {
		return &utils_errors.ValidationError{Field: "Email", Message: "campo obrigatório"}
	}
	if len(u.Email) > 100 {
		return &utils_errors.ValidationError{Field: "Email", Message: "máximo de 100 caracteres"}
	}
	if !utils_validators.IsValidEmail(u.Email) {
		return &utils_errors.ValidationError{Field: "Email", Message: "email inválido"}
	}

	// Password obrigatório, mínimo 8 caracteres, complexidade mínima
	if utils_validators.IsBlank(u.Password) {
		return &utils_errors.ValidationError{Field: "Password", Message: "campo obrigatório"}
	}
	if len(u.Password) < 8 {
		return &utils_errors.ValidationError{Field: "Password", Message: "mínimo de 8 caracteres"}
	}

	if !hasPasswordComplexity(u.Password) {
		return &utils_errors.ValidationError{Field: "Password", Message: "deve conter letras maiúsculas, minúsculas e números"}
	}

	return nil
}

// hasPasswordComplexity verifica se a senha tem pelo menos uma letra maiúscula,
// uma minúscula e um número.
func hasPasswordComplexity(pwd string) bool {
	var hasUpper, hasLower, hasNumber bool
	for _, c := range pwd {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}
	return hasUpper && hasLower && hasNumber
}
