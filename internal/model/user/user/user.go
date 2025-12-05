package model

import (
	"strings"
	"time"

	"unicode"

	valContact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type User struct {
	UID         int64
	Username    string
	Email       string
	Password    string
	Description string
	Status      bool
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate valida para criação (password obrigatório)
func (u *User) Validate() error {
	var errs []validators.ValidationError

	// Username
	if validators.IsBlank(u.Username) {
		errs = append(errs, validators.ValidationError{Field: "Username", Message: "campo obrigatório"})
	} else if len(u.Username) < 3 || len(u.Username) > 50 {
		errs = append(errs, validators.ValidationError{Field: "Username", Message: "deve ter entre 3 e 50 caracteres"})
	}

	// Email
	if validators.IsBlank(u.Email) {
		errs = append(errs, validators.ValidationError{Field: "Email", Message: "campo obrigatório"})
	} else if len(u.Email) > 100 {
		errs = append(errs, validators.ValidationError{Field: "Email", Message: "máximo de 100 caracteres"})
	} else if !valContact.IsValidEmail(u.Email) {
		errs = append(errs, validators.ValidationError{Field: "Email", Message: "email inválido"})
	}

	// Password (para criação)
	if validators.IsBlank(u.Password) {
		errs = append(errs, validators.ValidationError{Field: "Password", Message: "campo obrigatório"})
	} else if len(u.Password) < 8 {
		errs = append(errs, validators.ValidationError{Field: "Password", Message: "mínimo de 8 caracteres"})
	} else if !hasPasswordComplexity(u.Password) {
		errs = append(errs, validators.ValidationError{Field: "Password", Message: "deve conter letras maiúsculas, minúsculas e números"})
	}

	if len(errs) > 0 {
		return validators.NewValidationErrors(errs)
	}
	return nil
}

// ValidateForUpdate valida para atualização (password opcional)
func (u *User) ValidateForUpdate() error {
	var errs []validators.ValidationError

	// Username (obrigatório sempre)
	if validators.IsBlank(u.Username) {
		errs = append(errs, validators.ValidationError{Field: "Username", Message: "campo obrigatório"})
	} else if len(u.Username) < 3 || len(u.Username) > 50 {
		errs = append(errs, validators.ValidationError{Field: "Username", Message: "deve ter entre 3 e 50 caracteres"})
	}

	// Email (obrigatório sempre)
	if validators.IsBlank(u.Email) {
		errs = append(errs, validators.ValidationError{Field: "Email", Message: "campo obrigatório"})
	} else if len(u.Email) > 100 {
		errs = append(errs, validators.ValidationError{Field: "Email", Message: "máximo de 100 caracteres"})
	} else if !valContact.IsValidEmail(u.Email) {
		errs = append(errs, validators.ValidationError{Field: "Email", Message: "email inválido"})
	}

	// Password (opcional para update - valida apenas se fornecido E não for hash)
	if u.Password != "" && !isBcryptHash(u.Password) {
		// Apenas valida se não parece ser um hash
		if len(u.Password) < 8 {
			errs = append(errs, validators.ValidationError{Field: "Password", Message: "mínimo de 8 caracteres"})
		} else if !hasPasswordComplexity(u.Password) {
			errs = append(errs, validators.ValidationError{Field: "Password", Message: "deve conter letras maiúsculas, minúsculas e números"})
		}
	}

	if len(errs) > 0 {
		return validators.NewValidationErrors(errs)
	}
	return nil
}

func isBcryptHash(s string) bool {
	return strings.HasPrefix(s, "$2a$") ||
		strings.HasPrefix(s, "$2b$") ||
		strings.HasPrefix(s, "$2y$")
}

// hasPasswordComplexity verifica complexidade da senha
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
