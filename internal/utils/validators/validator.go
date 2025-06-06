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

func IsValidPhone(phone string) bool {
	// Ex: (11) 1234-5678
	re := regexp.MustCompile(`^\(\d{2}\) \d{4}-\d{4}$`)
	return re.MatchString(phone)
}

func IsValidCell(cell string) bool {
	// Ex: (11) 91234-5678
	re := regexp.MustCompile(`^\(\d{2}\) 9\d{4}-\d{4}$`)
	return re.MatchString(cell)
}

func IsValidCPF(cpf string) bool {
	// Simplesmente verifica se tem 11 dígitos numéricos (sem pontuação)
	matched, _ := regexp.MatchString(`^\d{11}$`, cpf)
	return matched
}

func IsValidCNPJ(cnpj string) bool {
	// Simplesmente verifica se tem 14 dígitos numéricos (sem pontuação)
	matched, _ := regexp.MatchString(`^\d{14}$`, cnpj)
	return matched
}
