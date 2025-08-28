package validators

import "regexp"

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
