package model

import (
	"strings"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

const (
	DefaultLimit     = 50
	MaxLimit         = 1000
	DefaultSortOrder = "asc"
)

type BaseFilter struct {
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
	SearchTerm string
}

// NewBaseFilter cria uma nova instância de BaseFilter com valores padrão aplicados
func NewBaseFilter(limit, offset int, sortBy, sortOrder, searchTerm string) *BaseFilter {
	b := &BaseFilter{
		Limit:      limit,
		Offset:     offset,
		SortBy:     sortBy,
		SortOrder:  strings.ToLower(sortOrder),
		SearchTerm: searchTerm,
	}
	b.applyDefaults()
	return b
}

// applyDefaults aplica valores padrão aos campos
func (b *BaseFilter) applyDefaults() {
	if b.Limit <= 0 {
		b.Limit = DefaultLimit
	}
	if b.Limit > MaxLimit {
		b.Limit = MaxLimit
	}
	if b.Offset < 0 {
		b.Offset = 0
	}
	if b.SortOrder == "" {
		b.SortOrder = DefaultSortOrder
	}
}

func (b *BaseFilter) Validate() error {
	if b.Limit < 0 {
		return &validators.ValidationError{Field: "Limit", Message: "não pode ser negativo"}
	}
	if b.Limit > MaxLimit {
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

// WithDefaults retorna uma cópia com valores padrão aplicados
func (b *BaseFilter) WithDefaults() BaseFilter {
	copy := *b
	copy.applyDefaults()
	return copy
}

// IsEmpty verifica se todos os campos de filtro estão vazios
func (b *BaseFilter) IsEmpty() bool {
	return b.Limit == 0 &&
		b.Offset == 0 &&
		b.SortBy == "" &&
		b.SortOrder == "" &&
		b.SearchTerm == ""
}
