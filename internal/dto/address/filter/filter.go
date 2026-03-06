package dto

import (
	"errors"
	"strings"
	"time"

	filterAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address/filter"
	commonFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
)

type AddressFilterDTO struct {
	UserID      *int64     `schema:"user_id"`
	ClientCpfID *int64     `schema:"client_cpf_id"`
	SupplierID  *int64     `schema:"supplier_id"`
	City        string     `schema:"city"`
	State       string     `schema:"state"`
	PostalCode  string     `schema:"postal_code"`
	IsActive    *bool      `schema:"is_active"`
	CreatedFrom *time.Time `schema:"created_from"`
	CreatedTo   *time.Time `schema:"created_to"`
	UpdatedFrom *time.Time `schema:"updated_from"`
	UpdatedTo   *time.Time `schema:"updated_to"`
	Limit       int        `schema:"limit"`
	Offset      int        `schema:"offset"`
	SortBy      string     `schema:"sort_by"`
	SortOrder   string     `schema:"sort_order"`
}

func (d *AddressFilterDTO) ToModel() (*filterAddress.AddressFilter, error) {
	filter := &filterAddress.AddressFilter{
		BaseFilter: commonFilter.BaseFilter{
			Limit:     d.Limit,
			Offset:    d.Offset,
			SortBy:    d.SortBy,
			SortOrder: d.SortOrder,
		},
		UserID:      d.UserID,
		ClientCpfID: d.ClientCpfID,
		SupplierID:  d.SupplierID,
		City:        d.City,
		State:       d.State,
		PostalCode:  d.PostalCode,
		IsActive:    d.IsActive,
		CreatedFrom: d.CreatedFrom,
		CreatedTo:   d.CreatedTo,
		UpdatedFrom: d.UpdatedFrom,
		UpdatedTo:   d.UpdatedTo,
	}

	return filter, nil
}

func (d *AddressFilterDTO) Validate() error {
	var validationErrors []string

	city := strings.TrimSpace(d.City)
	state := strings.TrimSpace(d.State)
	postalCode := strings.TrimSpace(d.PostalCode)

	// Pelo menos um filtro de conteúdo
	hasContentFilter :=
		d.UserID != nil ||
			d.ClientCpfID != nil ||
			d.SupplierID != nil ||
			city != "" ||
			state != "" ||
			postalCode != "" ||
			d.IsActive != nil ||
			d.CreatedFrom != nil ||
			d.CreatedTo != nil ||
			d.UpdatedFrom != nil ||
			d.UpdatedTo != nil

	if !hasContentFilter {
		validationErrors = append(validationErrors, "pelo menos um filtro de busca deve ser fornecido")
	}

	// ===== VALIDAÇÕES ESPECÍFICAS =====
	if city != "" && len(city) < 2 {
		validationErrors = append(validationErrors, "'city' deve conter no mínimo 2 caracteres")
	}

	if state != "" && len(state) != 2 {
		validationErrors = append(validationErrors, "'state' deve conter exatamente 2 caracteres (UF)")
	}

	if postalCode != "" {
		cleanPostalCode := strings.ReplaceAll(postalCode, "-", "")
		cleanPostalCode = strings.ReplaceAll(cleanPostalCode, ".", "")
		if len(cleanPostalCode) != 8 {
			validationErrors = append(validationErrors, "'postal_code' inválido - deve conter 8 dígitos")
		}
	}

	// ===== PAGINAÇÃO =====
	if d.Limit < 1 {
		validationErrors = append(validationErrors, "'limit' deve ser maior que zero")
	}
	if d.Limit > 100 {
		validationErrors = append(validationErrors, "'limit' não pode ser maior que 100")
	}
	if d.Offset < 0 {
		validationErrors = append(validationErrors, "'offset' não pode ser negativo")
	}
	if d.Offset > 10_000 {
		validationErrors = append(validationErrors, "'offset' excede o limite permitido")
	}

	// ===== ORDENAÇÃO =====
	if d.SortBy != "" && !isValidAddressSortField(d.SortBy) {
		validationErrors = append(validationErrors, "'sort_by' inválido")
	}

	if d.SortOrder != "" {
		order := strings.ToLower(d.SortOrder)
		if order != "asc" && order != "desc" {
			validationErrors = append(validationErrors, "'sort_order' inválido")
		}
	}

	// ===== DATAS =====
	if d.CreatedFrom != nil && d.CreatedTo != nil && d.CreatedFrom.After(*d.CreatedTo) {
		validationErrors = append(validationErrors, "'created_from' não pode ser maior que 'created_to'")
	}

	if d.UpdatedFrom != nil && d.UpdatedTo != nil && d.UpdatedFrom.After(*d.UpdatedTo) {
		validationErrors = append(validationErrors, "'updated_from' não pode ser maior que 'updated_to'")
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}
	return nil
}

func isValidAddressSortField(field string) bool {
	allowedFields := map[string]bool{
		"id":            true,
		"user_id":       true,
		"client_cpf_id": true,
		"supplier_id":   true,
		"city":          true,
		"state":         true,
		"postal_code":   true,
		"is_active":     true,
		"created_at":    true,
		"updated_at":    true,
	}
	return allowedFields[strings.ToLower(field)]
}
