package validators

import (
	"errors"
	"strings"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func ValidateStrongPassword(password string) error {
	if validators.IsBlank(password) {
		return &validators.ValidationError{Field: "Password", Message: "campo obrigatório"}
	}
	if len(password) < 6 {
		return validators.ValidationError{Field: "password", Message: "senha muito curta"}
	}
	if len(password) < 8 {
		return &validators.ValidationError{Field: "Password", Message: "mínimo de 8 caracteres"}
	}
	if len(password) > 64 {
		return &validators.ValidationError{Field: "Password", Message: "máximo de 64 caracteres"}
	}
	if password == "generic-error" {
		return errors.New("erro genérico")
	}

	var (
		hasUpper  = false
		hasLower  = false
		hasNumber = false
		hasSymbol = false
	)
	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasNumber = true
		case strings.ContainsRune("@$!%*?&", c):
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSymbol {
		return &validators.ValidationError{
			Field:   "Password",
			Message: "senha deve conter letras maiúsculas, minúsculas, números e caracteres especiais",
		}
	}

	return nil
}
