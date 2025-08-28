package validators

import "regexp"

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`

	re := regexp.MustCompile(regex)
	return re.MatchString(email)
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
