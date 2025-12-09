package model

import (
	"strings"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type BaseFilter struct {
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
	SearchTerm string
}

func (b *BaseFilter) Validate() error {
	if b.Limit < 0 {
		return &validators.ValidationError{Field: "Limit", Message: "não pode ser negativo"}
	}
	if b.Limit > 1000 {
		return &validators.ValidationError{Field: "Limit", Message: "máximo permitido é 1000"}
	}
	if b.Offset < 0 {
		return &validators.ValidationError{Field: "Offset", Message: "não pode ser negativo"}
	}
	b.SortOrder = strings.ToLower(b.SortOrder)

	if b.SortOrder != "" && b.SortOrder != "asc" && b.SortOrder != "desc" {
		return &validators.ValidationError{Field: "SortOrder", Message: "deve ser 'asc' ou 'desc'"}
	}
	return nil
}
