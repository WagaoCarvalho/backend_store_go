package validators

import (
	"regexp"
)

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
