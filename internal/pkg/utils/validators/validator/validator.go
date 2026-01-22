package validators

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("erro no campo '%s': %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	out := ""
	for _, err := range ve {
		out += err.Error() + "; "
	}
	return out
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func EqualsIgnoreCaseAndTrim(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func ValidateSingleNonNil(fields ...*int64) bool {
	count := 0
	for _, f := range fields {
		if f != nil {
			count++
		}
	}
	return count == 1
}

// NewValidationErrors cria um novo erro de validação a partir de uma lista de ValidationError
func NewValidationErrors(errors []ValidationError) error {
	if len(errors) == 0 {
		return nil
	}
	return ValidationErrors(errors)
}

func IsEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// validação estrutural simples e segura (RFC-compatível para APIs)
	if len(email) > 254 {
		return false
	}

	at := strings.LastIndex(email, "@")
	if at < 1 || at == len(email)-1 {
		return false
	}

	local := email[:at]
	domain := email[at+1:]

	if len(local) > 64 {
		return false
	}

	// domínio precisa ter pelo menos um ponto
	if !strings.Contains(domain, ".") {
		return false
	}

	// caracteres inválidos óbvios
	if strings.ContainsAny(email, " <>(),;:\\\"[]") {
		return false
	}

	return true
}
