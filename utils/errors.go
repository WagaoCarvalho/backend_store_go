package utils

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("registro não encontrado")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Erro no campo '%s': %s", e.Field, e.Message)
}
