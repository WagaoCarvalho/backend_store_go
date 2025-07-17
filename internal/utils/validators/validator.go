package utils

import (
	"regexp"
	"strings"
)

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`

	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

var postalCodeRegex = regexp.MustCompile(`^\d{5}-?\d{3}$`)

func IsValidPostalCode(code string) bool {
	if code == "00000-000" || code == "00000000" {
		return false
	}
	return postalCodeRegex.MatchString(code)
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

	matched, _ := regexp.MatchString(`^\d{11}$`, cpf)
	return matched
}

func IsValidCNPJ(cnpj string) bool {

	matched, _ := regexp.MatchString(`^\d{14}$`, cnpj)
	return matched
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
