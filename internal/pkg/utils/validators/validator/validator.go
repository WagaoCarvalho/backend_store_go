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
	return strings.TrimSpace(s) == ""
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
