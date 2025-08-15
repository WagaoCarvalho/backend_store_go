package utils

import (
	"regexp"
	"strings"

	"github.com/WagaoCarvalho/backend_store_go/pkg/utils"
)

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`

	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

var postalCodeRegex = regexp.MustCompile(`^\d{8}$`)

func IsValidPostalCode(code string) bool {
	if !postalCodeRegex.MatchString(code) {
		return false
	}

	firstChar := code[0]
	for i := 1; i < len(code); i++ {
		if code[i] != firstChar {
			return true
		}
	}

	return false
}

var validStates = map[string]bool{
	"AC": true, "AL": true, "AP": true, "AM": true, "BA": true,
	"CE": true, "DF": true, "ES": true, "GO": true, "MA": true,
	"MT": true, "MS": true, "MG": true, "PA": true, "PB": true,
	"PR": true, "PE": true, "PI": true, "RJ": true, "RN": true,
	"RS": true, "RO": true, "RR": true, "SC": true, "SP": true,
	"SE": true, "TO": true,
}

func IsValidBrazilianState(state string) bool {
	return validStates[state]
}

func NormalizePhone(phone string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(phone, "")
}

func IsValidPhone(input string) bool {
	normalized := NormalizePhone(input)
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(normalized)
}

func IsValidCell(input string) bool {
	normalized := NormalizePhone(input)
	re := regexp.MustCompile(`^\d{11}$`)
	return re.MatchString(normalized)
}

func IsValidCPF(cpf string) bool {
	if match := regexp.MustCompile(`^\d{11}$`).MatchString(cpf); !match {
		return false
	}

	firstChar := cpf[0]
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != firstChar {
			return true
		}
	}
	return false
}

func IsValidCNPJ(cnpj string) bool {
	if match := regexp.MustCompile(`^\d{14}$`).MatchString(cnpj); !match {
		return false
	}

	firstChar := cnpj[0]
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != firstChar {
			return true
		}
	}
	return false
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

func ValidateStrongPassword(password string) error {
	if IsBlank(password) {
		return &utils.ValidationError{Field: "Password", Message: "campo obrigatório"}
	}
	if len(password) < 8 {
		return &utils.ValidationError{Field: "Password", Message: "mínimo de 8 caracteres"}
	}
	if len(password) > 64 {
		return &utils.ValidationError{Field: "Password", Message: "máximo de 64 caracteres"}
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
		return &utils.ValidationError{
			Field:   "Password",
			Message: "senha deve conter letras maiúsculas, minúsculas, números e caracteres especiais",
		}
	}

	return nil
}
