package utils

import "regexp"

func IsStringEmpty(param string) bool {
	if param == "" {
		return true
	}
	return false
}

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
